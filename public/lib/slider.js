/**
 * Created by tsned on 3/23/17.
 */
'use strict';

const _slider_tmpl = '' +
    '<div class="track-container">' +
    '  <div class="track-slider" layout="row" layout-align="center center">' +
    '    <span class="time-label">{{track.current | trackTime}}</span>' +
    '    <div class="track-bar">' +
    '      <div class="track-progress" ng-style="progressStyle"></div>' +
    '    </div>' +
    '    <span class="time-label">{{track.total | trackTime}}</span>' +
    '  </div>' +
    '</div>';

const _slider = angular.module('gotunes.track.slider', ['ngMaterial']);

_slider.factory('nSliderFactory', [function () {
    return function (scope, elem) {
        const self = this;
        self.scope = scope;
        self.elem = elem;
        self.scope.progressStyle = { width: '0%' };
        self.scope.track = { total: 0, current: 0 };
        self.barElem = self.elem.find('.track-bar');

        self.updateProgress = function (curr, tot) {
            let prog = (curr / (tot || 1)) * 100;
            self.scope.progressStyle.width = (prog > 100 ? 0 : prog) + '%';
        };

        self.getEventPosPercent = function (event) {
            let barWidth = this.barElem.width();
            let barOffset = this.barElem.offset().left;
            return ((event.clientX - barOffset) / (barWidth || 1)) * 100;
        };

        self.scope.$watch('current', function (newVal, oldVal) {
            if (newVal !== oldVal) {
                self.scope.track.current = newVal;
                self.updateProgress(self.scope.track.current, self.scope.track.total);
            }
        });

        self.scope.$watch('total', function (newVal, oldVal) {
            if (newVal !== oldVal) {
                self.scope.track.total = newVal;
                self.updateProgress(self.scope.track.current, self.scope.track.total);
            }
        });

        self.scope.$on('$destroy', function () {
            self.unbindEvents();
        });

        self.bindEvents = function () {
            self.barElem.on('mousedown', angular.bind(self, self.onSeek));
            //self.barElem.on('touchstart', angular.bind(self, self.onSeek));
        };

        self.unbindEvents = function () {
            self.barElem.off();
        };

        self.onSeek = function (event) {
            event.stopPropagation();
            event.preventDefault();
            let newPos = this.getEventPosPercent(event);
            if (this.scope.onSeek)
                this.scope.onSeek(newPos);
        };

        // do the bind
        self.bindEvents();
    };
}]);

_slider.directive('nTrackSlider', ['nSliderFactory', function (Slider) {
    return {
        restrict: 'E',
        replace: true,
        scope: {
            current: '=',
            total: '=',
            onSeek: '=?'
        },
        template: _slider_tmpl,
        link: function (scope, elem) {
            scope.slider = new Slider(scope, elem);
        }
    }
}]);

_slider.filter('trackTime', function () {
    return function (val) {
        let dur = moment.duration(val * 1000);
        return moment(dur.asMilliseconds()).format('mm:ss');
    }
});