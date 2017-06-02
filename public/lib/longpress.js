/**
 * Created by tsned on 5/25/17.
 */
(function (angular) {
    'use strict';
    let lp = angular.module('ngLongPress', []);
    lp.directive('ngLongPress', ['$parse', '$timeout', function ($parse, $timeout) {
        return {
            restrict: 'A',
            scope: true,
            link: function (scope, elem, attrs) {
                let delay = attrs.delay || 500;
                let noop = function () { return angular.noop };
                let longFn = attrs.ngLongPress ? $parse(attrs.ngLongPress) : noop;
                let shortFn = attrs.shortPress ? $parse(attrs.shortPress) : noop;
                let timeout = null;

                elem.css({'-webkit-user-select': 'none'});

                let mouseup = function (event) {
                    $timeout.cancel(timeout);
                    shortFn(scope, {$event: event});
                };

                elem.on('mousedown', function (event) {
                    event.preventDefault();
                    event.stopPropagation();

                    elem.one('mouseup', mouseup);

                    timeout = $timeout(function () {
                        elem.off('mouseup', mouseup);
                        longFn(scope, {$event: event});
                    }, delay);

                    return false;
                });
            }
        }
    }]);
})(window.angular);