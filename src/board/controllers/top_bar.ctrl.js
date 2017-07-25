(function(angular) {
    'use strict';

    angular.module('gitlabKBApp.board').controller('TopBarController', [
        '$scope',
        '$state',
        '$stateParams',
        'BoardService',
        'AuthService',
        '$window',
        'project_id',
        'project_path',
        'board_full_view',
        function($scope, $state, $stateParams, BoardService, AuthService, $window, project_id, project_path, board_full_view) {

			var projectId = $stateParams.project_id || project_id;
			var projectPath = $stateParams.project_path || project_path;

            if (projectPath !== undefined) {
                BoardService.get(projectPath).then(function(board) {
                    $scope.project = board.project;
                });
            }

            $scope.stateParams = $stateParams;

			$scope.topBarClass = board_full_view ? "fixed": "";

            $scope.logout = function() {
                AuthService.logout();
                $window.location.pathname = '/';
            };

            if (board_full_view) {
				$scope.switchViewLabel = "Return";
			} else  {
				$scope.switchViewLabel = "Full Screen";
			}

			$scope.switchView = function() {
				//AuthService.logout();
				if (board_full_view) {
					$window.location.pathname = '/' + project_path + "/board/";
				} else  {
					$window.location.pathname = '/' + project_path + "/board/full/";
				}
			};

            $scope.showActionBar = function() {
                BoardService.get(projectPath).then(function(board) {
                    board.state.showActionBar = !board.state.showActionBar;
                });
            };

            $scope.reset = function() {
                var params = {
                    project_id: projectId,
                    project_path: projectPath,
                    group: ''
                };
                $state.go('board.cards', params);
            };

            $scope.group = function(field) {
                var params = {
                    project_id: projectId,
                    project_path: projectPath,
                    group: field
                };
                $state.go('board.cards', params);
            };
        }
    ]);
}(window.angular));
