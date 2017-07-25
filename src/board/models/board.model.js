(function (angular) {
    'use strict';

    angular.module('gitlabKBApp.board').factory('Board',
        [
            'UserService',
            'Stage',
            'State',
            'LabelService',
            '$rootScope',
            function (UserService, Stage, State, LabelService, $rootScope) {
                function Board(labels, issues, project) {
                    this.stages = [];

                    this.issues = [];
                    this.stale = false;
                    this.project = project;
                    this.grouped = false;
                    this.defaultStages = {};
                    this.state = new State();
                    this.counts = {};
                    this.stages = LabelService.listStages(project.path_with_namespace);
                    this.priorities = LabelService.listPriorities(project.path_with_namespace);
                    this.viewLabels = LabelService.listViewLabels(project.path_with_namespace);
                    this.priorityLabels = _.map(this.priorities, 'name');
                    this.stagelabels = _.map(this.stages, 'name');
                    _.each(this.stages, _.bind(function (stage) {
                        this.defaultStages[stage.viewName] = [];
                    }, this));

                    this.initViewLabels = function (issue) {
                        issue.viewLabels = [];

                        var issueStage = _.find(issue.labels, function(l){return l.group == 'stage';});
                        if (issueStage) {
							issue.stage = LabelService.getStage(project.path_with_namespace, issueStage.name );
						}
						//issue.stage = this.getCardStage(issue);

                        if (! issue.stage) {
                           // issue.stage = this.stages[0];
                            issue.stage = new Stage();
                        }

						var issuePriority = _.find(issue.labels, function(l){return l.group == 'priority';});
						if (issuePriority) {
							issue.priority = LabelService.getPriority(project.path_with_namespace, issuePriority.name );
						} else {
							issue.priority = LabelService.getPriority(project.path_with_namespace, "");
						}

                       // issue.priority = LabelService.getPriority(project.path_with_namespace, _.intersection(this.priorityLabels, issue.labels)[0]);
						//issue.priority = _.find(issue.labels, function(l){return l.group == 'priority'});

						// issue.priority = LabelService.getStage(project.path_with_namespace,
						// 	_.find(issue.labels, function(l){return l.group == 'priority';}
						// 	).name);

                        if (!_.isEmpty(issue.labels)) {
                            var labels = issue.labels;
                            for (var i = 0; i < labels.length; i++) {
                                var label = this.viewLabels[labels[i].name];
                                if (label !== undefined) {
                                    issue.viewLabels.push(label);
                                }
                            }
                        }

                        return issue;
                    };

                    this.issues = _.map(issues, _.bind(this.initViewLabels, this));

                    this.byStage = function (element, index, items) {
                        element = _.chain(element)
                                  .sortBy(function(item) {return item.id * -1})
                                  .sortBy(function(item) {return item.priority.index * -1})
                                  .value();

                        var stages = {};
                        for (var k in this.defaultStages) {
                            stages[k] = [];
                        }

                        items[index] = _.extend(stages, _.groupBy(element, function(el){
                            return el.stage.viewName;
                        }, this));
                        for (var idx in items[index]) {
                            this.counts[idx] += items[index][idx].length;
                        }
                    };

                    this.add = function (card) {
                        var old = this.getCardById(card.id);

                        if (_.isEmpty(old)) {
                            // card.stage = LabelService.getStage(project.path_with_namespace,
                            //     _.find(card.labels, function(l){return l.group == 'stage'}
                            // ));
							card.stage = this.getCardStage(card);
                            this.initViewLabels(card);
                            this.issues.push(card);
                            $rootScope.$emit('board.change');
                        }
                    };

                    this.getCardStage = function (issue) {
						var issueStage = _.find(issue.labels, function(l){return l.group == 'stage';});
						if (issueStage) {
							return LabelService.getStage(project.path_with_namespace, issueStage.name );
						}
						return null;
					};

                    this.update = function (card) {
                        var old = this.getCardById(card.id);
                        _.extend(old, card);
                        // old.stage = LabelService.getStage(project.path_with_namespace,
                        //     _.find(old.labels, function(l){return l.group == 'stage'}
                        // ));
                        old.stage = this.getCardStage(card);
                        this.initViewLabels(old);
                        $rootScope.$emit('board.change');
                    };

                    this.remove = function (card) {
                        var old = this.getCardById(card.id);
                        this.issues.splice(this.issues.indexOf(old), 1);
                        $rootScope.$emit('board.change');
                    };

                    this.getCardById = function (id) {
                        return _.find(this.issues, function (item) {
                            return item.id == id;
                        });
                    };

                    this.reset = function(filter, group) {
                        for (var k in this.defaultStages) {
                            this.counts[k] = 0;
                        }

                        var issues = _.filter(this.issues, filter);
                        var groups = _.groupBy(issues, group);
                        groups = _.isEmpty(groups) ? {
                            none: {}
                        } : groups;
                        groups = _.each(groups, _.bind(this.byStage, this));

                        return groups;
                    };

                    this.listCard = function(filter) {
                        return _.chain(this.issues)
                                .filter(filter)
                                .sortBy(function(item) {return item.id * -1})
                                .sortBy(function(item) {return item.priority.index * -1})
                                .value()
                    }
                }

                return Board;
            }
        ]
    );
})(window.angular);
