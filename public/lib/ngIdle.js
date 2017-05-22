(function (angular) {
    'use strict';

    const idle = angular.module('ngIdle', []);

    idle.service('ngIdle', ['$document', '$interval', '$rootScope', function ($document, $interval, $rootScope) {
        const self = this;
        let count = 0;
        let idled = false;
        let interval = null;
        const ignored = [];
        const events = 'keydown keypress mousemove mousedown click wheel select touchstart touchmove';
        const options = {
            checkInterval: 2, // check every 2 seconds
            idle: 30 // 30 seconds until idle
        };

        const ignorer = function (e) {
            e.stopPropagation();
            return true;
        };

        self.setCheckInterval = function (s) { options.checkInterval = Math.abs(s) };
        self.getCheckInterval = function () { return options.checkInterval };
        self.setIdle = function (s) { options.idle = Math.abs(s) };
        self.getIdle = function () { return options.idle };
        self.reset = function () { count = 0 };

        self.stop = function () {
            count = 0;
            idled = false;
            $interval.cancel(interval);
            ignored.forEach(function (i) { self.unignore(i) });
            ignored.length = 0;
        };

        self.start = function () {
            if (!interval) {
                interval = $interval(function () {
                    count += options.checkInterval;
                    if (count > options.idle && !idled) {
                        $rootScope.$broadcast('ngIdle');
                        idled = true;
                    }
                }, options.checkInterval * 1000);
            }
        };

        self.ignore = function (elem) {
            elem.on(events, ignorer);
            ignored.push(elem);
        };

        self.unignore = function (elem) {
            elem.off(events, ignorer);
            let ix = ignored.indexOf(elem);
            if (ix > -1) elem.splice(ix, 1);
        };

        $document.on(events, function () {
            count = 0;
            idled = false;
        });
    }]);
})(window.angular);
