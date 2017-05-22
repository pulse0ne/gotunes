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
            const RequestViewType = ((window.enums || {}).Command || {}).REQUEST_VIEW;
            const ContextType = (window.enums || {}).ContextType || {};
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

            //$scope.$watch('nowplaying.track', function (nv, ov) {
            //    if (!nv) {
            //        $scope.showStatusBar = false;
            //    } else if (nv !== ov) {
            //        let track = $scope.nowplaying.track;
            //        let artist = track.artist || 'Unknown';
            //        let song = track.title || track.filename || 'Unknown';
            //        $scope.statusBarText = '"' + song + '" by ' + artist;
            //        $scope.showStatusBar = true;
            //    }
            //}, true);

            $scope.changeView = function (view) {
                switch(view) {
                    case ContextType.ALL_ARTISTS:
                    case ContextType.ALL_ALBUMS:
                    case ContextType.ALL_TRACKS:
                    case ContextType.PLAYLIST:
                    case ContextType.SCAN_ROOT:
                        gotunes.sendCommand(RequestViewType, view);
                        break;
                    default:
                        console.error('Unrecognized top-level view: ' + view);
                        break;
                }
            };

            $scope.playInContext = function (track) {

            };
        }
    ]);
})(window.angular);