(function(angular) {
    'use strict';
    angular.module('gitlabKBApp.board').factory('LabelService', [
        '$q',
        '$http',
        'CardPriority',
        'Stage',
        function($q, $http, CardPriority, Stage) {
            return {
                labels: [],
                list: function(path_with_namespace, withCache) {
                    withCache = (typeof withCache === 'undefined') ? true : withCache;
                    return $http.get('/api/boards/' + path_with_namespace + "/labels", {
                        cache: withCache
                    }).then(function(result) {
                        var labels = result.data.data;
                        if (_.isEmpty(labels)) {
                            return {};
                        }

                        var reserved = _.sortBy(_.filter(labels, function(label) {
                            return  label.group == 'stage';
                        }), 'name');

                        if (_.isEmpty(reserved)) {
                            return {};
                        }

                        this.labels[path_with_namespace] = labels;
                        return this.labels[path_with_namespace];
                    }.bind(this));
                },
                listStages: function(projectId){
                    return _.chain(this.labels[projectId])
                        .filter(function(label) {
                            return label.group == 'stage';
                        })
                        .map(function(label){
                            return  new Stage(label);
                        })
                        .sortBy(function(label){
                            return label.index;
                        }).value();
                },
                listPriorities: function(projectId) {
                    return _.chain(this.labels[projectId])
                        .filter(function(label) {
                            return label.group == 'priority' ;
                        })
                        .map(function(label){
                            return  new CardPriority(label);
                        })
                        .sortBy(function(label){
                            return label.index * -1;
                        }).value();
                },
                listViewLabels: function(projectId) {
                    return _.chain(this.labels[projectId])
                           .filter(function(label) {
                                return !(label.group == 'priority' || label.group == 'stage');
                           })
                           .keyBy('name')
                           .value();
                },
                getPriority: function(projectId, label){
                    var priority =_.find(this.labels[projectId], {name: label});
                    return new CardPriority(priority);
                },
                getStage: function(projectId, label) {
                    if (_.isEmpty(label)) {
                        return "";
                    }
                    var foundedStage = new Stage({name: label});
                    var stage = _.find(this.listStages(projectId), {viewName: foundedStage.viewName});

                    if (_.isEmpty(stage)) {
                        return "";
                    }

                    return new Stage(stage);
                },
                getStageByName: function(projectId, viewName) {
                    return this.list(projectId).then(function(labels){
                        var stages = _.chain(labels)
                            .filter(function(label) {
                                return label.group == 'stage';
                            }).map(function(label){
                                return  new Stage(label);
                            }).value();

                        var stage = _.find(stages, {viewName: viewName});

                        if (_.isEmpty(stage)) {
                            return null;
                        }

                        return stage;
                    });
                },
                create: function(projectId, label, color) {
                    return $http.post("/api/labels/" + projectId, {
                        name: label,
                        color: color
                    });
                },
                update: function(projectId, oldLabel, newLabel, color) {
                    return $http.put('/api/labels/' + projectId, {
                        name: oldLabel,
                        color: color,
                        new_name: newLabel
                    });
                },
                delete: function(projectId, label) {
                    return $http.delete("/api/labels/" + projectId + "/" + label);
                }
            };
        }
    ]);
})(window.angular);
