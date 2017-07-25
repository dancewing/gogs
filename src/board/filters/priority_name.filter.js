(function (angular) {
    'use strict';

    angular.module('gitlabKBApp.board').filter('priorityName', [
        function () {
            return function (input) {
                // var priority = input.match(priority_regexp);
                // if (_.isEmpty(priority)) {
                //     return input;
                // }
                // return priority[2];
				console.log(input);
				return false;
            };
        }]);
}(window.angular));
