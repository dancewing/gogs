(function(angular, _) {
    'use strict';

    angular.module('gitlabKBApp.board').controller('FilterController', [
        '$scope',
        '$state',
        '$stateParams',
        'BoardService',
        'UserService',
        'MilestoneService',
        'project_id',
        'project_path',
        function($scope, $state, $stateParams, BoardService, UserService, MilestoneService, project_id, project_path) {
            var labels = [],
                milestones = [],
                users = [],
                priority = [];

            var projectId = $stateParams.project_id || project_id;
            var projectPath = $stateParams.project_path || project_path;

            this.tags = _.isArray($stateParams.tags) ? $stateParams.tags : [$stateParams.tags];

            BoardService.get(projectPath).then(function(board) {
                this.labels = _.toArray(board.viewLabels);
                this.priorities = board.priorities;

                MilestoneService.list(board.project.id, board.project.path_with_namespace).then(function(milestones) {
                    this.milestones = milestones;
                }.bind(this));

                UserService.list(board.project.id, board.project.path_with_namespace).then(function(users) {
                    this.users = users;
                }.bind(this));

            }.bind(this));

            /**
             * Apply selected filtering criteria
             */
            this.apply = function(tag) {
                var params = {
                    project_id: projectId,
                    project_path: projectPath,
                    tags: $stateParams.tags
                };

                if (_.isArray(params.tags)) {
                    var idx = params.tags.indexOf(tag);

                    if (idx == -1) {
                        params.tags = params.tags.concat([tag]);
                    } else {
                        params.tags.splice(idx, 1);
                    }
                } else if (_.isString(params.tags)) {
                    if (params.tags == tag) {
                        params.tags = [];
                    } else {
                        params.tags = [params.tags, tag];
                    }
                } else {
                    params.tags = [tag];
                }

                $state.go('board.cards', params);
            };

            this.applyAll = function(tagPrefix, tags, identifyBy, enable) {
                var params = {
                    project_id: projectId,
                    project_path: projectPath,
                    tags: Array.isArray($stateParams.tags) ? $stateParams.tags : $stateParams.tags ? [$stateParams.tags] : []
                };

                tags = _(tags).values().map(function(tag) { return tagPrefix + tag[identifyBy]; }).value().concat(tagPrefix);

                if (enable) {
                    params.tags = _.uniq((params.tags || []).concat(tags));
                } else {
                    params.tags && _.pullAll(params.tags, tags);
                }

                $state.go('board.cards', params);
            };

            /**
             * Clear all filters
             */
            this.clear = function() {
                $state.go('board.cards', {
                    tags: []
                });
            };

            this.checked = function(obj) {
                return _.includes(this.tags, obj);
            }
        }
    ]);
}(window.angular, window._));
