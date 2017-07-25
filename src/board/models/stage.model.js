(function(angular) {
    'use strict';

    angular.module('gitlabKBApp.board').factory('Stage',[function() {
            function Stage(label) {
                if (_.isEmpty(label)) {
                    return {
                        index: 0,
                        color: '#fffff'
                    }
                }

				return {
					id: label.id,
					name: label.name,
					index: label.order,
					color: label.color,
					viewName: label.name,
					wip: undefined
				};
            }

            return Stage;
        }]
    );
})(window.angular);
