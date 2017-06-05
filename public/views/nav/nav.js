(function (angular) {
    'use strict';
    const app = angular.module('gotunes');

    app.controller('gotunes.controller.nav', [
        '$scope',
        '$document',
        '$mdMedia',
        'gotunes.service',
        function ($scope, $document, $mdMedia, gotunes) {
            $scope.viewClass = 'nav';
            $scope.$mdMedia = $mdMedia;
            $scope.sidenavWidth = { width: 0 };
            const ViewType = (window.enums || {}).ViewType || {};
            const Command = (window.enums || {}).Command || {};
            const RequestViewType = Command.REQUEST_VIEW;
            const tlviews = [
                ViewType.QUEUE,
                ViewType.ALL_ARTISTS,
                ViewType.ALL_ALBUMS,
                ViewType.ALL_TRACKS,
                ViewType.PLAYLIST
            ];
            let sidenav = angular.element('#left-sidenav');
            let navspacer = angular.element('#sidenav-spacer');
            let header = angular.element('#content-header');
            let headspacer = angular.element('#content-header-spacer');

            // keep sidenav spacer width in sync
            $scope.$watch(
                function () { return sidenav.outerWidth() },
                function (nv) { navspacer.width(nv) }
            );

            // keep header spacer height in sync
            $scope.$watch(
                function () { return header.outerHeight() },
                function (nv) { headspacer.height(nv) }
            );

            $scope.changeView = function (view) {
                if (tlviews.includes(view)) {
                    gotunes.sendCommand(RequestViewType, view);
                } else {
                    console.error('Unrecognized top-level view:', view);
                }
            };

            $scope.addToQueue = function (track) {
                gotunes.sendCommand(Command.ADD_TO_QUEUE, track.file);
            };

            $scope.openItemMenu = function (item) {
                console.log('open item menu');
                $scope.selected === item ? delete $scope.selected : $scope.selected = item;
            };

            $scope.testcap = {};
        }
    ]);
})(window.angular);