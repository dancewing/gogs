(function(angular) {
    'use strict';

    angular.module('gitlabKBApp.board').factory('CommentService',
        [
            '$http','project_path', function($http, project_path) {
                return {
                    list: function(boardId, cardId) {
                        return $http.get('/api/boards/'+project_path+'/cards/'+ cardId +'/comments', {
                            params: {
                                project_id: boardId,
                                issue_id: cardId
                            }
                        }).then(function(result) {
                            return result.data.data;
                        });
                    },
                    create: function(boardId, cardId, comment) {
                        return $http.post('/api/boards/'+project_path+'/cards/'+ cardId +'/comments', {
                            project_id: boardId,
                            issue_id: cardId,
                            body: comment
                        }).then(function(result) {
                            return result.data.data;
                        });
                    }
                };
            }
        ]
    );
})(window.angular);

