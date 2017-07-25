(function(angular){
    'use strict';

    angular.module('gitlabKBApp.board').factory('CardPriority', [
        function() {
            function CardPriority(label) {
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
					viewName: label.name
				};

            }

            return CardPriority;
        }
    ]);

}(window.angular));
