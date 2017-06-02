(function (angular) {
    'use strict';
    const app = angular.module('gotunes');

    app.controller('gotunes.controller.idle', [
        '$scope',
        '$rootScope',
        '$document',
        '$location',
        '$interval',
        '$timeout',
        'fitTextService',
        function ($scope, $rootScope, $document, $location, $interval, $timeout, fitTextService) {
            $scope.viewClass = 'idle';

            let animationInterval;
            let trackWatchTimeout;
            let wallIx = Math.floor(Math.random() * 5);
            const animationDuration = 60000;
            const events = 'mousedown touchstart keydown';
            const walls = [1, 2, 3, 4, 5].map(function (w) {
                return 'assets/img/wall' + w + '.jpg'
            });

            const handler = function (evt) {
                evt.preventDefault();
                evt.stopPropagation();
                $rootScope.showStatusBar = true;
                $location.path('/');
            };

            const getNextImageUrl = function () {
                return ++wallIx >= walls.length ? walls[wallIx = 0] : walls[wallIx];
            };

            const switchImage = function () {
                let next = getNextImageUrl();
                let el = angular.element('#wall');
                el.removeClass('idle-animation');
                el.attr('src', next);
                el.addClass('idle-animation');
            };

            $document.on(events, handler);

            $scope.$on('$destroy', function () {
                $document.off(events, handler);
                $interval.cancel(animationInterval);
                $timeout.cancel(trackWatchTimeout);
            });

            animationInterval = $interval(switchImage, animationDuration);

            $scope.$watch('nowplaying.track', function (nv, ov) {
                if (nv !== ov) {
                    trackWatchTimeout = $timeout(function () {
                        fitTextService.forceRecalculate('artist title');
                    });
                }
            }, true);

            switchImage();
        }
    ]);
})(window.angular);