/**
 * Created by tsned on 6/2/17.
 */
(function (angular) {
    'use strict';
    let m = angular.module('ngKeybindings', []);
    m.directive('keyCapture', ['$parse', '$document', '$keybindings', function ($parse, $document, $keybindings) {
        return {
            restrict: 'E',
            replace: true,
            scope: { keyModel: '=', autoBind: '=' },
            template:
                '<div class="key-capture-container" style="display: flex; flex-direction: row;">' +
                '  <button class="key-capture-activate-button" style="margin-right: 10px;" ng-class="{\'capture-active\': active}">{{active ? "Capturing..." : "New Capture"}}</button>' +
                '  <div class="key-presses" style="display: flex; flex-direction: row;">' +
                '    <div class="kp-tag-wrapper" ng-repeat="k in keyModel.keys">' +
                '      <span class="kp-tag" style="border: 1px solid rgba(0,0,0,0.5); border-radius: 10%; margin: 4px; padding: 4px;">{{k}}</span>' +
                '      <span class="kp-plus" ng-if="!$last && !($first && $last)">+</span>' +
                '    </div>' +
                '  </div>' +
                '</div>',
            link: function (scope, elem) {
                let keyCodes = {
                    3 : 'break',
                    8 : 'backspace',
                    9 : 'tab',
                    12 : 'clear',
                    13 : 'enter',
                    16 : 'shift',
                    17 : 'ctrl',
                    18 : 'alt',
                    19 : 'pause/break',
                    20 : 'caps lock',
                    27 : 'escape',
                    32 : 'spacebar',
                    33 : 'page up',
                    34 : 'page down',
                    35 : 'end',
                    36 : 'home',
                    37 : 'left arrow',
                    38 : 'up arrow',
                    39 : 'right arrow',
                    40 : 'down arrow',
                    41 : 'select',
                    42 : 'print',
                    43 : 'execute',
                    44 : 'print screen',
                    45 : 'insert',
                    46 : 'delete',
                    48 : '0',
                    49 : '1',
                    50 : '2',
                    51 : '3',
                    52 : '4',
                    53 : '5',
                    54 : '6',
                    55 : '7',
                    56 : '8',
                    57 : '9',
                    58 : ':',
                    60 : '<',
                    61 : '=',
                    64 : '@',
                    65 : 'a',
                    66 : 'b',
                    67 : 'c',
                    68 : 'd',
                    69 : 'e',
                    70 : 'f',
                    71 : 'g',
                    72 : 'h',
                    73 : 'i',
                    74 : 'j',
                    75 : 'k',
                    76 : 'l',
                    77 : 'm',
                    78 : 'n',
                    79 : 'o',
                    80 : 'p',
                    81 : 'q',
                    82 : 'r',
                    83 : 's',
                    84 : 't',
                    85 : 'u',
                    86 : 'v',
                    87 : 'w',
                    88 : 'x',
                    89 : 'y',
                    90 : 'z',
                    91 : 'meta',
                    96 : '0',
                    97 : '1',
                    98 : '2',
                    99 : '3',
                    100 : '4',
                    101 : '5',
                    102 : '6',
                    103 : '7',
                    104 : '8',
                    105 : '9',
                    106 : '*',
                    107 : '+',
                    108 : '.',
                    109 : '-',
                    110 : '.',
                    111 : '/',
                    112 : 'f1',
                    113 : 'f2',
                    114 : 'f3',
                    115 : 'f4',
                    116 : 'f5',
                    117 : 'f6',
                    118 : 'f7',
                    119 : 'f8',
                    120 : 'f9',
                    121 : 'f10',
                    122 : 'f11',
                    123 : 'f12',
                    124 : 'f13',
                    125 : 'f14',
                    126 : 'f15',
                    127 : 'f16',
                    128 : 'f17',
                    129 : 'f18',
                    130 : 'f19',
                    131 : 'f20',
                    132 : 'f21',
                    133 : 'f22',
                    134 : 'f23',
                    135 : 'f24',
                    144 : 'num lock',
                    145 : 'scroll lock',
                    160 : '^',
                    161: '!',
                    163 : '#',
                    164: '$',
                    166 : 'page backward',
                    167 : 'page forward',
                    170: '*',
                    186 : ';',
                    187 : '=',
                    188 : ',',
                    189 : '-',
                    190 : '.',
                    191 : '/',
                    219 : '[',
                    220 : '\\',
                    221 : ']',
                    222 : '\'',
                    223 : '`',
                };

                let mods = [16, 17, 18, 91];
                let button = angular.element(elem.children()[0]);

                scope.active = false;
                if (!scope.keyModel) {
                    scope.keyModel = {
                        text: '',
                        code: 0,
                        isAlt: false,
                        isCtrl: false,
                        isMeta: false,
                        isShift: false
                    };
                }
                if (!scope.keyModel.keys) scope.keyModel.keys = [];

                let keydownFn = function (evt) {
                    evt.preventDefault();
                    evt.stopPropagation();

                    if (mods.includes(evt.keyCode)) {
                        if (!scope.keyModel.keys.includes(keyCodes[evt.keyCode])) {
                            scope.keyModel.keys.push(keyCodes[evt.keyCode]);
                        }
                        if (evt.keyCode === 16) scope.keyModel.isShift = true;
                        else if (evt.keyCode === 17) scope.keyModel.isCtrl = true;
                        else if (evt.keyCode === 18) scope.keyModel.isAlt = true;
                        else if (evt.keyCode === 91) scope.keyModel.isMeta = true;
                    } else {
                        scope.keyModel.keys.push(keyCodes[evt.keyCode]);
                        scope.keyModel.code = evt.keyCode;
                        scope.keyModel.text = keyCodes[evt.keyCode];
                    }
                    scope.$apply();
                };

                let keyupFn = function (evt) {
                    evt.preventDefault();
                    evt.stopPropagation();

                    if (mods.includes(evt.keyCode) && scope.active) {
                        let ix = scope.keyModel.keys.indexOf(keyCodes[evt.keyCode]);
                        if (ix > -1) scope.keyModel.keys.splice(ix, 1);
                    } else {
                        scope.active = false;
                    }
                    scope.$apply();
                };

                let resetModel = function () {
                    scope.keyModel.keys.length = 0;
                    for (let k in scope.keyModel) {
                        if (scope.keyModel.hasOwnProperty(k) && k.startsWith('is')) {
                            scope.keyModel[k] = false;
                        }
                    }
                    scope.keyModel.text = '';
                    scope.keyModel.code = 0;
                };

                scope.$watch('active', function (nv, ov) {
                    if (nv) {
                        resetModel();
                        $keybindings.setEnabled(false);
                        $document.on('keydown', keydownFn);
                        $document.on('keyup', keyupFn);
                    } else {
                        if (scope.keyModel.keys.length === 0) {
                            resetModel();
                        }
                        $document.off('keydown', keydownFn);
                        $document.off('keyup', keyupFn);
                        $keybindings.setEnabled(true);
                    }
                });

                button.on('click', function () {
                    scope.active = !scope.active;
                    scope.$apply();
                });
            }
        };
    }]);

    m.service('$keybindings', ['$document', function ($document) {
        let self = this;
        let enabled = false;
        let bindings = {};

        let genId = function (binding) {
            if (!!binding.code) {
                return '_' + binding.code + '_' +
                    ((binding.isAlt << 1) +
                    (binding.isCtrl << 2) +
                    (binding.isMeta << 3) +
                    (binding.isShift << 4)) + '_';
            } else {
                return '_' + binding.keyCode + '_' +
                    ((binding.altKey << 1) +
                    (binding.ctrlKey << 2) +
                    (binding.metaKey << 3) +
                    (binding.shiftKey << 4)) + '_';
            }
        };

        self.checkEvent = function (binding, evt) {
            return binding.code === evt.keyCode &&
                binding.isAlt === evt.altKey &&
                binding.isCtrl === evt.ctrlKey &&
                binding.isMeta === evt.metaKey &&
                binding.isShift === evt.shiftKey;
        };

        self.setEnabled = function (bool) {
            enabled = bool;
        };

        self.addBinding = function (binding, func) {
            bindings[genId(binding)] = func;
        };

        self.removeBinding = function (binding) {
            delete bindings[genId(binding)];
        };

        $document.on('keydown', function (evt) {
            if (!enabled) {
                return true;
            }

            let binding = bindings[genId(evt)];
            if (!binding) {
                console.log('no binding');
                return true;
            }

            binding(evt);
        });

    }]);
})(window.angular);