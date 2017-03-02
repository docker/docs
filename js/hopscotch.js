/**! hopscotch - v0.2.6
 *
 * Copyright 2016 LinkedIn Corp. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
! function(a, b) {
    "use strict";
    if ("function" == typeof define && define.amd) define([], b);
    else if ("object" == typeof exports) module.exports = b();
    else {
        var c = "hopscotch";
        if (a[c]) return;
        a[c] = b()
    }
}(this, function() {
    var Hopscotch, HopscotchBubble, HopscotchCalloutManager, HopscotchI18N, customI18N, customRenderer, customEscape, utils, callbacks, helpers, winLoadHandler, defaultOpts, winHopscotch, templateToUse = "bubble_default",
        Sizzle = window.Sizzle || null,
        undefinedStr = "undefined",
        waitingToStart = !1,
        hasJquery = typeof jQuery !== undefinedStr,
        hasSessionStorage = !1,
        isStorageWritable = !1,
        document = window.document,
        validIdRegEx = /^[a-zA-Z]+[a-zA-Z0-9_-]*$/,
        rtlMatches = {
            left: "right",
            right: "left"
        };
    try {
        typeof window.sessionStorage !== undefinedStr && (hasSessionStorage = !0, sessionStorage.setItem("hopscotch.test.storage", "ok"), sessionStorage.removeItem("hopscotch.test.storage"), isStorageWritable = !0)
    } catch (err) {}
    return defaultOpts = {
            smoothScroll: !0,
            scrollDuration: 1e3,
            scrollTopMargin: 200,
            showCloseButton: !0,
            showPrevButton: !1,
            showNextButton: !0,
            bubbleWidth: 280,
            bubblePadding: 15,
            arrowWidth: 20,
            skipIfNoElement: !0,
            isRtl: !1,
            cookieName: "hopscotch.tour.state"
        }, Array.isArray || (Array.isArray = function(a) {
            return "[object Array]" === Object.prototype.toString.call(a)
        }), winLoadHandler = function() {
            waitingToStart && winHopscotch.startTour()
        }, utils = {
            addClass: function(a, b) {
                var c, d, e, f;
                if (a.className) {
                    for (d = b.split(/\s+/), c = " " + a.className + " ", e = 0, f = d.length; f > e; ++e) c.indexOf(" " + d[e] + " ") < 0 && (c += d[e] + " ");
                    a.className = c.replace(/^\s+|\s+$/g, "")
                } else a.className = b
            },
            removeClass: function(a, b) {
                var c, d, e, f;
                for (d = b.split(/\s+/), c = " " + a.className + " ", e = 0, f = d.length; f > e; ++e) c = c.replace(" " + d[e] + " ", " ");
                a.className = c.replace(/^\s+|\s+$/g, "")
            },
            hasClass: function(a, b) {
                var c;
                return a.className ? (c = " " + a.className + " ", -1 !== c.indexOf(" " + b + " ")) : !1
            },
            getPixelValue: function(a) {
                var b = typeof a;
                return "number" === b ? a : "string" === b ? parseInt(a, 10) : 0
            },
            valOrDefault: function(a, b) {
                return typeof a !== undefinedStr ? a : b
            },
            invokeCallbackArrayHelper: function(a) {
                var b;
                return Array.isArray(a) && (b = helpers[a[0]], "function" == typeof b) ? b.apply(this, a.slice(1)) : void 0
            },
            invokeCallbackArray: function(a) {
                var b, c;
                if (Array.isArray(a)) {
                    if ("string" == typeof a[0]) return utils.invokeCallbackArrayHelper(a);
                    for (b = 0, c = a.length; c > b; ++b) utils.invokeCallback(a[b])
                }
            },
            invokeCallback: function(a) {
                return "function" == typeof a ? a() : "string" == typeof a && helpers[a] ? helpers[a]() : utils.invokeCallbackArray(a)
            },
            invokeEventCallbacks: function(a, b) {
                var c, d, e = callbacks[a];
                if (b) return this.invokeCallback(b);
                for (c = 0, d = e.length; d > c; ++c) this.invokeCallback(e[c].cb)
            },
            getScrollTop: function() {
                var a;
                return a = typeof window.pageYOffset !== undefinedStr ? window.pageYOffset : document.documentElement.scrollTop
            },
            getScrollLeft: function() {
                var a;
                return a = typeof window.pageXOffset !== undefinedStr ? window.pageXOffset : document.documentElement.scrollLeft
            },
            getWindowHeight: function() {
                return window.innerHeight || document.documentElement.clientHeight
            },
            addEvtListener: function(a, b, c) {
                return a ? a.addEventListener ? a.addEventListener(b, c, !1) : a.attachEvent("on" + b, c) : void 0
            },
            removeEvtListener: function(a, b, c) {
                return a ? a.removeEventListener ? a.removeEventListener(b, c, !1) : a.detachEvent("on" + b, c) : void 0
            },
            documentIsReady: function() {
                return "complete" === document.readyState
            },
            evtPreventDefault: function(a) {
                a.preventDefault ? a.preventDefault() : event && (event.returnValue = !1)
            },
            extend: function(a, b) {
                var c;
                for (c in b) b.hasOwnProperty(c) && (a[c] = b[c])
            },
            getStepTargetHelper: function(a) {
                var b = document.getElementById(a);
                if (b) return b;
                if (hasJquery) return b = jQuery(a), b.length ? b[0] : null;
                if (Sizzle) return b = new Sizzle(a), b.length ? b[0] : null;
                if (document.querySelector) try {
                    return document.querySelector(a)
                } catch (c) {}
                return /^#[a-zA-Z][\w-_:.]*$/.test(a) ? document.getElementById(a.substring(1)) : null
            },
            getStepTarget: function(a) {
                var b;
                if (!a || !a.target) return null;
                if ("string" == typeof a.target) return utils.getStepTargetHelper(a.target);
                if (Array.isArray(a.target)) {
                    var c, d;
                    for (c = 0, d = a.target.length; d > c; c++)
                        if ("string" == typeof a.target[c] && (b = utils.getStepTargetHelper(a.target[c]))) return b;
                    return null
                }
                return a.target
            },
            getI18NString: function(a) {
                return customI18N[a] || HopscotchI18N[a]
            },
            setState: function(a, b, c) {
                var d, e = "";
                if (hasSessionStorage && isStorageWritable) try {
                    sessionStorage.setItem(a, b)
                } catch (f) {
                    isStorageWritable = !1, this.setState(a, b, c)
                } else hasSessionStorage && sessionStorage.removeItem(a), c && (d = new Date, d.setTime(d.getTime() + 24 * c * 60 * 60 * 1e3), e = "; expires=" + d.toGMTString()), document.cookie = a + "=" + b + e + "; path=/"
            },
            getState: function(a) {
                var b, c, d, e = a + "=",
                    f = document.cookie.split(";");
                if (hasSessionStorage && (d = sessionStorage.getItem(a))) return d;
                for (b = 0; b < f.length; b++) {
                    for (c = f[b];
                        " " === c.charAt(0);) c = c.substring(1, c.length);
                    if (0 === c.indexOf(e)) {
                        d = c.substring(e.length, c.length);
                        break
                    }
                }
                return d
            },
            clearState: function(a) {
                hasSessionStorage ? sessionStorage.removeItem(a) : this.setState(a, "", -1)
            },
            normalizePlacement: function(a) {
                !a.placement && a.orientation && (a.placement = a.orientation)
            },
            flipPlacement: function(a) {
                if (a.isRtl && !a._isFlipped) {
                    var b, c, d = ["orientation", "placement"];
                    a.xOffset && (a.xOffset = -1 * this.getPixelValue(a.xOffset));
                    for (c in d) b = d[c], a.hasOwnProperty(b) && rtlMatches.hasOwnProperty(a[b]) && (a[b] = rtlMatches[a[b]]);
                    a._isFlipped = !0
                }
            }
        }, utils.addEvtListener(window, "load", winLoadHandler), callbacks = {
            next: [],
            prev: [],
            start: [],
            end: [],
            show: [],
            error: [],
            close: []
        }, helpers = {}, HopscotchI18N = {
            stepNums: null,
            nextBtn: "Next",
            prevBtn: "Back",
            doneBtn: "Done",
            skipBtn: "Skip",
            closeTooltip: "&#215;"
        }, customI18N = {}, HopscotchBubble = function(a) {
            this.init(a)
        }, HopscotchBubble.prototype = {
            isShowing: !1,
            currStep: void 0,
            setPosition: function(a) {
                var b, c, d, e, f, g, h, i = utils.getStepTarget(a),
                    j = this.element,
                    k = this.arrowEl,
                    l = a.isRtl ? "right" : "left";
                if (utils.flipPlacement(a), utils.normalizePlacement(a), c = j.offsetWidth, b = j.offsetHeight, utils.removeClass(j, "fade-in-down fade-in-up fade-in-left fade-in-right"), d = i.getBoundingClientRect(), h = a.isRtl ? d.right - c : d.left, "top" === a.placement) e = d.top - b - this.opt.arrowWidth, f = h;
                else if ("bottom" === a.placement) e = d.bottom + this.opt.arrowWidth, f = h;
                else if ("left" === a.placement) e = d.top, f = d.left - c - this.opt.arrowWidth;
                else {
                    if ("right" !== a.placement) throw new Error("Bubble placement failed because step.placement is invalid or undefined!");
                    e = d.top, f = d.right + this.opt.arrowWidth
                }
                g = "center" !== a.arrowOffset ? utils.getPixelValue(a.arrowOffset) : a.arrowOffset, g ? "top" === a.placement || "bottom" === a.placement ? (k.style.top = "", "center" === g ? k.style[l] = Math.floor(c / 2 - k.offsetWidth / 2) + "px" : k.style[l] = g + "px") : ("left" === a.placement || "right" === a.placement) && (k.style[l] = "", "center" === g ? k.style.top = Math.floor(b / 2 - k.offsetHeight / 2) + "px" : k.style.top = g + "px") : (k.style.top = "", k.style[l] = ""), "center" === a.xOffset ? f = d.left + i.offsetWidth / 2 - c / 2 : f += utils.getPixelValue(a.xOffset), "center" === a.yOffset ? e = d.top + i.offsetHeight / 2 - b / 2 : e += utils.getPixelValue(a.yOffset), a.fixedElement || (e += utils.getScrollTop(), f += utils.getScrollLeft()), j.style.position = a.fixedElement ? "fixed" : "absolute", j.style.top = e + "px", j.style.left = f + "px"
            },
            render: function(a, b, c) {
                var d, e, f, g, h, i, j, k, l, m, n = this.element;
                if (a ? this.currStep = a : this.currStep && (a = this.currStep), this.opt.isTourBubble ? (g = winHopscotch.getCurrTour(), g && (e = g.customData, d = g.customRenderer, a.isRtl = a.hasOwnProperty("isRtl") ? a.isRtl : g.hasOwnProperty("isRtl") ? g.isRtl : this.opt.isRtl, f = g.unsafe, Array.isArray(g.steps) && (h = g.steps.length, i = this._getStepI18nNum(this._getStepNum(h - 1)), k = this._getStepNum(b) === this._getStepNum(h - 1)))) : (e = a.customData, d = a.customRenderer, f = a.unsafe, a.isRtl = a.hasOwnProperty("isRtl") ? a.isRtl : this.opt.isRtl), j = k ? utils.getI18NString("doneBtn") : a.showSkip ? utils.getI18NString("skipBtn") : utils.getI18NString("nextBtn"), utils.flipPlacement(a), utils.normalizePlacement(a), this.placement = a.placement, m = {
                        i18n: {
                            prevBtn: utils.getI18NString("prevBtn"),
                            nextBtn: j,
                            closeTooltip: utils.getI18NString("closeTooltip"),
                            stepNum: this._getStepI18nNum(this._getStepNum(b)),
                            numSteps: i
                        },
                        buttons: {
                            showPrev: utils.valOrDefault(a.showPrevButton, this.opt.showPrevButton) && this._getStepNum(b) > 0,
                            showNext: utils.valOrDefault(a.showNextButton, this.opt.showNextButton),
                            showCTA: utils.valOrDefault(a.showCTAButton && a.ctaLabel, !1),
                            ctaLabel: a.ctaLabel,
                            showClose: utils.valOrDefault(this.opt.showCloseButton, !0)
                        },
                        step: {
                            num: b,
                            isLast: utils.valOrDefault(k, !1),
                            title: a.title || "",
                            content: a.content || "",
                            isRtl: a.isRtl,
                            placement: a.placement,
                            padding: utils.valOrDefault(a.padding, this.opt.bubblePadding),
                            width: utils.getPixelValue(a.width) || this.opt.bubbleWidth,
                            customData: a.customData || {}
                        },
                        tour: {
                            isTour: this.opt.isTourBubble,
                            numSteps: h,
                            unsafe: utils.valOrDefault(f, !1),
                            customData: e || {}
                        }
                    }, "function" == typeof d) n.innerHTML = d(m);
                else if ("string" == typeof d) {
                    if (!winHopscotch.templates || "function" != typeof winHopscotch.templates[d]) throw new Error('Bubble rendering failed - template "' + d + '" is not a function.');
                    n.innerHTML = winHopscotch.templates[d](m)
                } else if (customRenderer) n.innerHTML = customRenderer(m);
                else {
                    if (!winHopscotch.templates || "function" != typeof winHopscotch.templates[templateToUse]) throw new Error('Bubble rendering failed - template "' + templateToUse + '" is not a function.');
                    n.innerHTML = winHopscotch.templates[templateToUse](m)
                }
                for (children = n.children, numChildren = children.length, l = 0; l < numChildren; l++) node = children[l], utils.hasClass(node, "hopscotch-arrow") && (this.arrowEl = node);
                return n.style.zIndex = "number" == typeof a.zindex ? a.zindex : "", this._setArrow(a.placement), this.hide(!1), this.setPosition(a), c && c(!a.fixedElement), this
            },
            _getStepNum: function(a) {
                var b, c, d = 0,
                    e = winHopscotch.getSkippedStepsIndexes(),
                    f = e.length;
                for (c = 0; f > c; c++) b = e[c], a > b && d++;
                return a - d
            },
            _getStepI18nNum: function(a) {
                var b = utils.getI18NString("stepNums");
                return b && a < b.length ? a = b[a] : a += 1, a
            },
            _setArrow: function(a) {
                utils.removeClass(this.arrowEl, "down up right left"), "top" === a ? utils.addClass(this.arrowEl, "down") : "bottom" === a ? utils.addClass(this.arrowEl, "up") : "left" === a ? utils.addClass(this.arrowEl, "right") : "right" === a && utils.addClass(this.arrowEl, "left")
            },
            _getArrowDirection: function() {
                return "top" === this.placement ? "down" : "bottom" === this.placement ? "up" : "left" === this.placement ? "right" : "right" === this.placement ? "left" : void 0
            },
            show: function() {
                var a = this,
                    b = "fade-in-" + this._getArrowDirection(),
                    c = 1e3;
                return utils.removeClass(this.element, "hide"), utils.addClass(this.element, b), setTimeout(function() {
                    utils.removeClass(a.element, "invisible")
                }, 50), setTimeout(function() {
                    utils.removeClass(a.element, b)
                }, c), this.isShowing = !0, this
            },
            hide: function(a) {
                var b = this.element;
                return a = utils.valOrDefault(a, !0), b.style.top = "", b.style.left = "", a ? (utils.addClass(b, "hide"), utils.removeClass(b, "invisible")) : (utils.removeClass(b, "hide"), utils.addClass(b, "invisible")), utils.removeClass(b, "animate fade-in-up fade-in-down fade-in-right fade-in-left"), this.isShowing = !1, this
            },
            destroy: function() {
                var a = this.element;
                a && a.parentNode.removeChild(a), utils.removeEvtListener(a, "click", this.clickCb)
            },
            _handleBubbleClick: function(a) {
                function b(c) {
                    return c === a.currentTarget ? null : utils.hasClass(c, "hopscotch-cta") ? "cta" : utils.hasClass(c, "hopscotch-next") ? "next" : utils.hasClass(c, "hopscotch-prev") ? "prev" : utils.hasClass(c, "hopscotch-close") ? "close" : b(c.parentElement)
                }
                var c;
                a = a || window.event;
                var d = a.target || a.srcElement;
                if (c = b(d), "cta" === c) this.opt.isTourBubble || winHopscotch.getCalloutManager().removeCallout(this.currStep.id), this.currStep.onCTA && utils.invokeCallback(this.currStep.onCTA);
                else if ("next" === c) winHopscotch.nextStep(!0);
                else if ("prev" === c) winHopscotch.prevStep(!0);
                else if ("close" === c) {
                    if (this.opt.isTourBubble) {
                        var e = winHopscotch.getCurrStepNum(),
                            f = winHopscotch.getCurrTour(),
                            g = e === f.steps.length - 1;
                        utils.invokeEventCallbacks("close"), winHopscotch.endTour(!0, g)
                    } else this.opt.onClose && utils.invokeCallback(this.opt.onClose), this.opt.id && !this.opt.isTourBubble ? winHopscotch.getCalloutManager().removeCallout(this.opt.id) : this.destroy();
                    utils.evtPreventDefault(a)
                }
            },
            init: function(a) {
                var b, c, d, e, f = document.createElement("div"),
                    g = this,
                    h = !1;
                this.element = f, e = {
                    showPrevButton: defaultOpts.showPrevButton,
                    showNextButton: defaultOpts.showNextButton,
                    bubbleWidth: defaultOpts.bubbleWidth,
                    bubblePadding: defaultOpts.bubblePadding,
                    arrowWidth: defaultOpts.arrowWidth,
                    isRtl: defaultOpts.isRtl,
                    showNumber: !0,
                    isTourBubble: !0
                }, a = typeof a === undefinedStr ? {} : a, utils.extend(e, a), this.opt = e, f.className = "hopscotch-bubble", e.isTourBubble ? (d = winHopscotch.getCurrTour(), d && utils.addClass(f, "tour-" + d.id)) : utils.addClass(f, "hopscotch-callout no-number"), b = function() {
                    !h && g.isShowing && (h = !0, setTimeout(function() {
                        g.setPosition(g.currStep), h = !1
                    }, 100))
                }, utils.addEvtListener(window, "resize", b), this.clickCb = function(a) {
                    g._handleBubbleClick(a)
                }, utils.addEvtListener(f, "click", this.clickCb), this.hide(), utils.documentIsReady() ? document.body.appendChild(f) : (document.addEventListener ? (c = function() {
                    document.removeEventListener("DOMContentLoaded", c), window.removeEventListener("load", c), document.body.appendChild(f)
                }, document.addEventListener("DOMContentLoaded", c, !1)) : (c = function() {
                    "complete" === document.readyState && (document.detachEvent("onreadystatechange", c), window.detachEvent("onload", c), document.body.appendChild(f))
                }, document.attachEvent("onreadystatechange", c)), utils.addEvtListener(window, "load", c))
            }
        }, HopscotchCalloutManager = function() {
            var a = {},
                b = {};
            this.createCallout = function(c) {
                var d;
                if (!c.id) throw new Error("Must specify a callout id.");
                if (!validIdRegEx.test(c.id)) throw new Error("Callout ID is using an invalid format. Use alphanumeric, underscores, and/or hyphens only. First character must be a letter.");
                if (a[c.id]) throw new Error("Callout by that id already exists. Please choose a unique id.");
                if (!utils.getStepTarget(c)) throw new Error("Must specify existing target element via 'target' option.");
                return c.showNextButton = c.showPrevButton = !1, c.isTourBubble = !1, d = new HopscotchBubble(c), a[c.id] = d, b[c.id] = c, d.render(c, null, function() {
                    d.show(), c.onShow && utils.invokeCallback(c.onShow)
                }), d
            }, this.getCallout = function(b) {
                return a[b]
            }, this.removeAllCallouts = function() {
                var b;
                for (b in a) a.hasOwnProperty(b) && this.removeCallout(b)
            }, this.removeCallout = function(c) {
                var d = a[c];
                a[c] = null, b[c] = null, d && d.destroy()
            }, this.refreshCalloutPositions = function() {
                var c, d, e;
                for (c in a) a.hasOwnProperty(c) && b.hasOwnProperty(c) && (d = a[c], e = b[c], d && e && d.setPosition(e))
            }
        }, Hopscotch = function(a) {
            var b, c, d, e, f, g, h, i, j = this,
                k = {},
                l = [],
                m = function(a) {
                    return b && b.element && b.element.parentNode || (b = new HopscotchBubble(d)), a && utils.extend(b.opt, {
                        bubblePadding: o("bubblePadding"),
                        bubbleWidth: o("bubbleWidth"),
                        showNextButton: o("showNextButton"),
                        showPrevButton: o("showPrevButton"),
                        showCloseButton: o("showCloseButton"),
                        arrowWidth: o("arrowWidth"),
                        isRtl: o("isRtl")
                    }), b
                },
                n = function() {
                    b && (b.destroy(), b = null)
                },
                o = function(a) {
                    return "undefined" == typeof d ? defaultOpts[a] : utils.valOrDefault(d[a], defaultOpts[a])
                },
                p = function() {
                    var a;
                    return a = !e || 0 > f || f >= e.steps.length ? null : e.steps[f]
                },
                q = function() {
                    j.nextStep()
                },
                r = function(a) {
                    var b, c, d, e, f, g, h = m(),
                        i = h.element,
                        j = utils.getPixelValue(i.style.top),
                        k = j + utils.getPixelValue(i.offsetHeight),
                        l = utils.getStepTarget(p()),
                        n = l.getBoundingClientRect(),
                        q = n.top + utils.getScrollTop(),
                        r = n.bottom + utils.getScrollTop(),
                        s = q > j ? j : q,
                        t = k > r ? k : r,
                        u = utils.getScrollTop(),
                        v = u + utils.getWindowHeight(),
                        w = s - o("scrollTopMargin");
                    s >= u && (s <= u + o("scrollTopMargin") || v >= t) ? a && a() : o("smoothScroll") ? typeof YAHOO !== undefinedStr && typeof YAHOO.env !== undefinedStr && typeof YAHOO.env.ua !== undefinedStr && typeof YAHOO.util !== undefinedStr && typeof YAHOO.util.Scroll !== undefinedStr ? (b = YAHOO.env.ua.webkit ? document.body : document.documentElement, d = YAHOO.util.Easing ? YAHOO.util.Easing.easeOut : void 0, c = new YAHOO.util.Scroll(b, {
                        scroll: {
                            to: [0, w]
                        }
                    }, o("scrollDuration") / 1e3, d), c.onComplete.subscribe(a), c.animate()) : hasJquery ? jQuery("body, html").animate({
                        scrollTop: w
                    }, o("scrollDuration"), a) : (0 > w && (w = 0), e = u > s ? -1 : 1, f = Math.abs(u - w) / (o("scrollDuration") / 10), (g = function() {
                        var b = utils.getScrollTop(),
                            c = b + e * f;
                        return e > 0 && c >= w || 0 > e && w >= c ? (c = w, a && a(), void window.scrollTo(0, c)) : (window.scrollTo(0, c), utils.getScrollTop() === b ? void(a && a()) : void setTimeout(g, 10))
                    })()) : (window.scrollTo(0, w), a && a())
                },
                s = function(a, b) {
                    var c, d, g;
                    f + a >= 0 && f + a < e.steps.length ? (f += a, d = p(), g = function() {
                        c = utils.getStepTarget(d), c ? (k[f] && delete k[f], b(f)) : (k[f] = !0, utils.invokeEventCallbacks("error"), s(a, b))
                    }, d.delay ? setTimeout(g, d.delay) : g()) : b(-1)
                },
                t = function(a, b) {
                    var c, d, g, h, i = m(),
                        j = this;
                    if (i.hide(), a = utils.valOrDefault(a, !0), c = p(), c.nextOnTargetClick && utils.removeEvtListener(utils.getStepTarget(c), "click", q), d = c, g = b > 0 ? d.multipage : f > 0 && e.steps[f - 1].multipage, h = function(c) {
                            var e;
                            if (-1 === c) return this.endTour(!0);
                            if (a && (e = b > 0 ? utils.invokeEventCallbacks("next", d.onNext) : utils.invokeEventCallbacks("prev", d.onPrev)), c === f) {
                                if (g) return void x();
                                e = utils.valOrDefault(e, !0), e ? this.showStep(c) : this.endTour(!1)
                            }
                        }, !g && o("skipIfNoElement")) s(b, function(a) {
                        h.call(j, a)
                    });
                    else if (f + b >= 0 && f + b < e.steps.length) {
                        if (f += b, c = p(), !utils.getStepTarget(c) && !g) return utils.invokeEventCallbacks("error"), this.endTour(!0, !1);
                        h.call(this, f)
                    } else if (f + b === e.steps.length) return this.endTour();
                    return this
                },
                u = function(a) {
                    var b, c, d, e = {};
                    for (b in a) a.hasOwnProperty(b) && "id" !== b && "steps" !== b && (e[b] = a[b]);
                    return i.call(this, e, !0), c = utils.getState(o("cookieName")), c && (d = c.split(":"), g = d[0], h = d[1], d.length > 2 && (l = d[2].split(",")), h = parseInt(h, 10)), this
                },
                v = function(a, b, c) {
                    var d, e;
                    if (f = a || 0, k = b || {}, d = p(), e = utils.getStepTarget(d)) return void c(f);
                    if (!e) {
                        if (utils.invokeEventCallbacks("error"), k[f] = !0, o("skipIfNoElement")) return void s(1, c);
                        f = -1, c(f)
                    }
                },
                w = function(a) {
                    function b() {
                        d.show(), utils.invokeEventCallbacks("show", c.onShow)
                    }
                    var c = e.steps[a],
                        d = m(),
                        g = utils.getStepTarget(c);
                    f !== a && p().nextOnTargetClick && utils.removeEvtListener(utils.getStepTarget(p()), "click", q), f = a, d.hide(!1), d.render(c, a, function(a) {
                        a ? r(b) : b(), c.nextOnTargetClick && utils.addEvtListener(g, "click", q)
                    }), x()
                },
                x = function() {
                    var a = e.id + ":" + f,
                        b = winHopscotch.getSkippedStepsIndexes();
                    b && b.length > 0 && (a += ":" + b.join(",")), utils.setState(o("cookieName"), a, 1)
                },
                y = function(a) {
                    a && this.configure(a)
                };
            this.getCalloutManager = function() {
                return typeof c === undefinedStr && (c = new HopscotchCalloutManager), c
            }, this.startTour = function(a, b) {
                var c, d, f = {},
                    i = this;
                if (!e) {
                    if (!a) throw new Error("Tour data is required for startTour.");
                    if (!a.id || !validIdRegEx.test(a.id)) throw new Error("Tour ID is using an invalid format. Use alphanumeric, underscores, and/or hyphens only. First character must be a letter.");
                    e = a, u.call(this, a)
                }
                if (typeof b !== undefinedStr) {
                    if (b >= e.steps.length) throw new Error("Specified step number out of bounds.");
                    d = b
                }
                if (!utils.documentIsReady()) return waitingToStart = !0, this;
                if ("undefined" == typeof d && e.id === g && typeof h !== undefinedStr) {
                    if (d = h, l.length > 0)
                        for (var j = 0, k = l.length; k > j; j++) f[l[j]] = !0
                } else d || (d = 0);
                return v(d, f, function(a) {
                    var b = -1 !== a && utils.getStepTarget(e.steps[a]);
                    return b ? (utils.invokeEventCallbacks("start"), c = m(), c.hide(!1), i.isActive = !0, void(utils.getStepTarget(p()) ? i.showStep(a) : (utils.invokeEventCallbacks("error"), o("skipIfNoElement") && i.nextStep(!1)))) : void i.endTour(!1, !1)
                }), this
            }, this.showStep = function(a) {
                var b = e.steps[a],
                    c = f;
                return utils.getStepTarget(b) ? (b.delay ? setTimeout(function() {
                    w(a)
                }, b.delay) : w(a), this) : (f = a, utils.invokeEventCallbacks("error"), void(f = c))
            }, this.prevStep = function(a) {
                return t.call(this, a, -1), this
            }, this.nextStep = function(a) {
                return t.call(this, a, 1), this
            }, this.endTour = function(a, b) {
                var c, d = m();
                return a = utils.valOrDefault(a, !0), b = utils.valOrDefault(b, !0), e && (c = p(), c && c.nextOnTargetClick && utils.removeEvtListener(utils.getStepTarget(c), "click", q)), f = 0, h = void 0, d.hide(), a && utils.clearState(o("cookieName")), this.isActive && (this.isActive = !1, e && b && utils.invokeEventCallbacks("end")), this.removeCallbacks(null, !0), this.resetDefaultOptions(), n(), e = null, this
            }, this.getCurrTour = function() {
                return e
            }, this.getCurrTarget = function() {
                return utils.getStepTarget(p())
            }, this.getCurrStepNum = function() {
                return f
            }, this.getSkippedStepsIndexes = function() {
                var a, b = [];
                for (a in k) b.push(a);
                return b
            }, this.refreshBubblePosition = function() {
                var a = p();
                return a && m().setPosition(a), this.getCalloutManager().refreshCalloutPositions(), this
            }, this.listen = function(a, b, c) {
                return a && callbacks[a].push({
                    cb: b,
                    fromTour: c
                }), this
            }, this.unlisten = function(a, b) {
                var c, d, e = callbacks[a];
                for (c = 0, d = e.length; d > c; ++c) e[c].cb === b && e.splice(c, 1);
                return this
            }, this.removeCallbacks = function(a, b) {
                var c, d, e, f;
                for (f in callbacks)
                    if (!a || a === f)
                        if (b)
                            for (c = callbacks[f], d = 0, e = c.length; e > d; ++d) c[d].fromTour && (c.splice(d--, 1), --e);
                        else callbacks[f] = [];
                return this
            }, this.registerHelper = function(a, b) {
                "string" == typeof a && "function" == typeof b && (helpers[a] = b)
            }, this.unregisterHelper = function(a) {
                helpers[a] = null
            }, this.invokeHelper = function(a) {
                var b, c, d = [];
                for (b = 1, c = arguments.length; c > b; ++b) d.push(arguments[b]);
                helpers[a] && helpers[a].call(null, d)
            }, this.setCookieName = function(a) {
                return d.cookieName = a, this
            }, this.resetDefaultOptions = function() {
                return d = {}, this
            }, this.resetDefaultI18N = function() {
                return customI18N = {}, this
            }, this.getState = function() {
                return utils.getState(o("cookieName"))
            }, i = function(a, b) {
                var c, e, f, g, h = ["next", "prev", "start", "end", "show", "error", "close"];
                for (d || this.resetDefaultOptions(), utils.extend(d, a), a && utils.extend(customI18N, a.i18n), f = 0, g = h.length; g > f; ++f) e = "on" + h[f].charAt(0).toUpperCase() + h[f].substring(1), a[e] && this.listen(h[f], a[e], b);
                return c = m(!0), this
            }, this.configure = function(a) {
                return i.call(this, a, !1)
            }, this.setRenderer = function(a) {
                var b = typeof a;
                return "string" === b ? (templateToUse = a, customRenderer = void 0) : "function" === b && (customRenderer = a), this
            }, this.setEscaper = function(a) {
                return "function" == typeof a && (customEscape = a), this
            }, y.call(this, a)
        }, winHopscotch = new Hopscotch,
        function() {
            var _ = {};
            _.escape = function(a) {
                return customEscape ? customEscape(a) : null == a ? "" : ("" + a).replace(new RegExp("[&<>\"']", "g"), function(a) {
                    return "&" == a ? "&amp;" : "<" == a ? "&lt;" : ">" == a ? "&gt;" : '"' == a ? "&quot;" : "'" == a ? "&#x27;" : void 0
                })
            }, this.templates = this.templates || {}, this.templates.bubble_default = function(obj) {
                function optEscape(a, b) {
                    return b ? _.escape(a) : a
                }
                obj || (obj = {});
                var __t, __p = "";
                _.escape, Array.prototype.join;
                with(obj) __p += '\n<div class="hopscotch-bubble-container" style="width: ' + (null == (__t = step.width) ? "" : __t) + "px; padding: " + (null == (__t = step.padding) ? "" : __t) + 'px;">\n  ', tour.isTour && (__p += '<span class="hopscotch-bubble-number">' + (null == (__t = i18n.stepNum) ? "" : __t) + "</span>"), __p += '\n  <div class="hopscotch-bubble-content">\n    ', "" !== step.title && (__p += '<h3 class="hopscotch-title">' + (null == (__t = optEscape(step.title, tour.unsafe)) ? "" : __t) + "</h3>"), __p += "\n    ", "" !== step.content && (__p += '<div class="hopscotch-content">' + (null == (__t = optEscape(step.content, tour.unsafe)) ? "" : __t) + "</div>"), __p += '\n  </div>\n  <div class="hopscotch-actions">\n    ', buttons.showPrev && (__p += '<button class="hopscotch-nav-button prev hopscotch-prev">' + (null == (__t = i18n.prevBtn) ? "" : __t) + "</button>"), __p += "\n    ", buttons.showCTA && (__p += '<button class="hopscotch-nav-button next hopscotch-cta">' + (null == (__t = buttons.ctaLabel) ? "" : __t) + "</button>"), __p += "\n    ", buttons.showNext && (__p += '<button class="hopscotch-nav-button next hopscotch-next">' + (null == (__t = i18n.nextBtn) ? "" : __t) + "</button>"), __p += "\n  </div>\n  ", buttons.showClose && (__p += '<button title="Close" class="hopscotch-bubble-close hopscotch-close">' + (null == (__t = i18n.closeTooltip) ? "" : __t) + "</button>"), __p += '\n</div>\n<div class="hopscotch-bubble-arrow-container hopscotch-arrow">\n  <div class="hopscotch-bubble-arrow-border"></div>\n  <div class="hopscotch-bubble-arrow"></div>\n</div>';
                return __p
            }
        }.call(winHopscotch), winHopscotch
});
