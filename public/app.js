/**
 * Created by tsned on 3/21/17.
 */
(function (angular) {
    'use strict';
    const app = angular.module('gotunes', [
        'ngRoute',
        'ngAnimate',
        'ngMaterial',
        'ngWebsocket',
        'ngIdle',
        'ngFitText',
        'gotunes.track.slider'
    ]);

    app.config(['$routeProvider', '$locationProvider', function ($routeProvider, $locationProvider) {
        $routeProvider.when('/', {
            templateUrl: 'views/nav/nav.html',
            controller: 'gotunes.controller.nav'
        });
        $routeProvider.when('/idle', {
            templateUrl: 'views/idle/idle.html',
            controller: 'gotunes.controller.idle'
        });
        $routeProvider.otherwise({
            redirectTo: '/'
        });

        $locationProvider.html5Mode(true);
    }]);

    app.constant('SVGs', {
        play: 'M 0 0 L 0 32 L 32 16 z',
        pause: 'M 2 0 L 2 32 L 12 32 L 12 0 L 2 0 M 20 0 L 20 32 L 30 32 L 30 0 L 20 0'
    });

    app.service('gotunes.service', [
        '$websocket',
        '$rootScope',
        function ($websocket, $rootScope) {
            const self = this;
            const MessageType = (window.enums || {}).MessageType || {};

            const ws = $websocket.$new({
                url: 'ws://' + window.location.host + '/websocket',
                reconnect: true,
                reconnectInterval: 1000,
                maxReconnectInterval: 10000,
                doubleInterval: true,
                immediate: false
            });

            self.sendCommand = function (cmdType, data) {
                ws.$send({type: MessageType.COMMAND, payload: {command: cmdType, data: data}});
            };

            ws.$on('$message', function (message) {
                message = JSON.parse(message);
                $rootScope.$broadcast('ws.message', message);
            });

            ws.$on('$open', function () {
                $rootScope.$broadcast('ws.open');
            });

            ws.$on('$close', function (code) {
                $rootScope.$broadcast('ws.close', code);
            });

            ws.$on('$error', function () {
                $rootScope.$broadcast('ws.error');
            });

            self.connect = ws.$open;
        }
    ]);

    app.controller('gotunes.controller.main', [
        '$scope',
        '$rootScope',
        '$timeout',
        '$location',
        '$document',
        '$mdDialog',
        'ngIdle',
        'gotunes.service',
        'SVGs',
        function ($scope, $rootScope, $timeout, $location, $document, $mdDialog, ngIdle, gotunes, svg) {
            const enums = window.enums || {};
            const Command = enums.Command || {};
            const PlayState = $scope.ps = enums.PlayState || {};
            const ContextType = $scope.context = enums.ContextType || {};
            const MessageType = enums.MessageType || {};
            const unicodePlay = '\u25B6 ';

            let lastVolume = 50;
            let domSlider;

            $scope.wsConnected = false;
            $scope.nowplaying = {
                playstate: PlayState.PAUSED,
                time: {
                    current: 0,
                    total: 0
                },
                track: null
            };
            $scope.view = {
                type: null,
                parent: null,
                data: []
            };
            $scope.svgPath = svg.play;

            // TODO: test code
            ngIdle.setIdle(5);
            // TODO

            $scope.$watch('nowplaying.volume', function (newVal) {
                if (newVal !== lastVolume) {
                    domSlider.noUiSlider.set(Math.max(Math.min(newVal, 100), 0));
                }
            });

            $scope.$watch('nowplaying.playstate', function (newVal) {
                $scope.svgPath = newVal !== PlayState.PLAYING ? svg.play : svg.pause;
                if (newVal === PlayState.PLAYING && $document[0].title.substring(0, unicodePlay.length) !== unicodePlay) {
                    $document[0].title = unicodePlay + $document[0].title;
                } else if ($document[0].title.substring(0, unicodePlay.length) === unicodePlay) {
                    $document[0].title = $document[0].title.substring(unicodePlay.length);
                }
            });

            $scope.$watch('nowplaying.track', function (nv, ov) {
                let track = $scope.nowplaying.track;
                if (!nv || (!track.artist && !track.title && !track.filename)) {
                    $document[0].title = 'gotunes';
                    $rootScope.showStatusBar = false;
                } else if (nv !== ov) {
                    let artist = track.artist || 'Unknown';
                    let song = track.title || track.filename || 'Unknown';
                    let ps = $document[0].title.substring(0, 1) === unicodePlay ? unicodePlay : '';
                    $document[0].title = ps + '"' + song + '" by ' + artist;
                    $rootScope.showStatusBar = true;
                    $rootScope.statusBarText = '"' + song + '" by ' + artist;
                }
            }, true);

            $scope.$on('ngIdle', function () {
                if ($location.path() !== '/idle' && $scope.nowplaying.playstate === PlayState.PLAYING) {
                    $rootScope.showStatusBar = false;
                    $location.path('/idle');
                }
            });

            $scope.$on('ws.open', function () {
                gotunes.sendCommand(Command.REQUEST_VIEW, ContextType.ALL_ARTISTS);
            });

            $scope.$on('ws.message', function (evt, msg) {
                switch(msg.type) {
                    case MessageType.VIEW_UPDATE:
                        $scope.view.type = msg.payload.type;
                        $scope.view.data = msg.payload.data;
                        break;
                    case MessageType.NOW_PLAYING:
                        $scope.nowplaying = msg.payload;
                        break;
                    default:
                        break;
                }
                $scope.$apply();
            });

            $scope.onSeek = function (percent) {
                gotunes.sendCommand(Command.SEEK_TO, percent);
            };

            $scope.togglePlaystate = function () {
                let newState = ($scope.nowplaying || {}).playstate === PlayState.PLAYING ? PlayState.PAUSED : PlayState.PLAYING;
                gotunes.sendCommand(Command.SET_PLAYSTATE, newState);
            };

            $scope.playNext = angular.bind(this, gotunes.sendCommand, Command.PLAY_NEXT);
            $scope.playPrevious = angular.bind(this, gotunes.sendCommand, Command.PLAY_PREV);

            domSlider = $document[0].getElementById('slider');
            noUiSlider.create(domSlider, {
                start: lastVolume,
                connect: [true, false],
                range: { min: 0, max: 100 }
            });

            domSlider.noUiSlider.on('start', function (val) {
                lastVolume = val;
            });

            domSlider.noUiSlider.on('end', function (val) {
                gotunes.sendCommand(Command.SET_VOLUME, Math.floor(val));
            });

            let ignored = angular.element('#sticky-footer');

            gotunes.connect();
            ngIdle.ignore(ignored);
            ngIdle.start();
        }
    ]);

    app.filter('trackTime', function () {
        return function (val) {
            let dur = moment.duration(val * 1000);
            let time = moment(dur.asMilliseconds()).format('mm:ss');
            return time === 'Invalid date' ? '00:00' : time;
        }
    });
})(window.angular);
