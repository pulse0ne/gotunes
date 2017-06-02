/**
 * Created by tsned on 5/26/17.
 */
(function (angular) {
    'use strict';
    let ls = angular.module('ngLocalStorage', []);
    ls.service('$localStorage', ['$window', '$rootScope', function ($window, $rootScope) {
        let self = this;

        self.save = function (key, value) {
            if (angular.isObject(value)) value = JSON.stringify(value);
            $window.localStorage.setItem(key, value);
            $rootScope.$broadcast('ngLocalStorage.save', key);
        };

        self.load = function (key) {
            let val = $window.localStorage.getItem(key);
            try {
                val = JSON.parse(val);
            } catch (e) {}
            $rootScope.$broadcast('ngLocalStorage.load', key, val);
            return val;
        };

        self.remove = function (key) {
            $window.localStorage.removeItem(key);
            $rootScope.$broadcast('ngLocalStorage.remove', key);
        };

        self.clear = function () {
            $window.localStorage.clear();
            $rootScope.$broadcast('ngLocalStorage.clear', key);
        };

        self.forEach = function (fn) {
            if (!fn || fn.constructor !== Function) return;
            let len = $window.localStorage.length;
            for (let i = 0; i < len; i++) {
                let k = $window.localStorage.key(i);
                fn(k, $window.localStorage.getItem(k));
            }
        };
    }]);
})(window.angular);