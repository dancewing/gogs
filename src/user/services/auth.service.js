(function(angular){
    'use strict';

    angular.module('gitlabKBApp.user').factory('AuthService',
        [
            '$http',
            '$q',
            'store',
            'project_path',
            function ($http, $q, store, project_path) {
                return {
                    current: undefined,
                    roles: {
                        anon: 0,
                        user: 1
                    },
                    authenticate: function (data) {
                        return $http.post('/api/login', {
                            _username: data.username,
                            _password: data.password
                        }).then(function (result) {
                            store.set('id_token', result.data.token);
                            return store.get('id_token');
                        });
                    },
                    getCurrent: function () {
                    	if (this.current === undefined) {
							return $http.get('/api/boards/'+ project_path +'/current').then(function (result) {
								this.current = result.data;
								return this.current;
							}.bind(this));
						}
						return $q.when(this.current);
                    },
                    isAuthenticated: function () {
                        return this.getCurrent() !== null;
                    },
                    authorized: function (state) {
                        var roles = this.roles;
                        return !!(this.isAuthenticated()
                        || state.data.access === undefined
                        || state.data.access == roles.anon);
                    }
                };
            }
        ]
    );
})(window.angular);


