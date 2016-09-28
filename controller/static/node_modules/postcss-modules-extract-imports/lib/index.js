'use strict';

Object.defineProperty(exports, '__esModule', {
  value: true
});

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { 'default': obj }; }

function _slicedToArray(arr, i) { if (Array.isArray(arr)) { return arr; } else if (Symbol.iterator in Object(arr)) { var _arr = []; var _n = true; var _d = false; var _e = undefined; try { for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) { _arr.push(_s.value); if (i && _arr.length === i) break; } } catch (err) { _d = true; _e = err; } finally { try { if (!_n && _i['return']) _i['return'](); } finally { if (_d) throw _e; } } return _arr; } else { throw new TypeError('Invalid attempt to destructure non-iterable instance'); } }

var _postcss = require('postcss');

var _postcss2 = _interopRequireDefault(_postcss);

var declWhitelist = ['composes'],
    declFilter = new RegExp('^(' + declWhitelist.join('|') + ')$'),
    matchImports = /^(.+?)\s+from\s+(?:"([^"]+)"|'([^']+)')$/,
    icssImport = /^:import\((?:"([^"]+)"|'([^']+)')\)/;

var processor = _postcss2['default'].plugin('modules-extract-imports', function (options) {
  return function (css) {
    var imports = {},
        importIndex = 0,
        createImportedName = options && options.createImportedName || function (importName /*, path*/) {
      return 'i__imported_' + importName.replace(/\W/g, '_') + '_' + importIndex++;
    };

    // Find any declaration that supports imports
    css.walkDecls(declFilter, function (decl) {
      var matches = decl.value.match(matchImports);
      if (matches) {
        (function () {
          var _matches = _slicedToArray(matches, 4);

          var symbols = _matches[1];
          var doubleQuotePath = _matches[2];
          var singleQuotePath = _matches[3];

          var path = doubleQuotePath || singleQuotePath;
          imports[path] = imports[path] || {};
          var tmpSymbols = symbols.split(/\s+/).map(function (s) {
            if (!imports[path][s]) {
              imports[path][s] = createImportedName(s, path);
            }
            return imports[path][s];
          });
          decl.value = tmpSymbols.join(' ');
        })();
      }
    });

    // If we've found any imports, insert or append :import rules
    var existingImports = {};
    css.walkRules(function (rule) {
      var matches = icssImport.exec(rule.selector);
      if (matches) {
        var _matches2 = _slicedToArray(matches, 3);

        var doubleQuotePath = _matches2[1];
        var singleQuotePath = _matches2[2];

        existingImports[doubleQuotePath || singleQuotePath] = rule;
      }
    });

    Object.keys(imports).reverse().forEach(function (path) {

      var rule = existingImports[path];
      if (!rule) {
        rule = _postcss2['default'].rule({
          selector: ':import("' + path + '")',
          raws: { after: '\n' }
        });
        css.prepend(rule);
      }
      Object.keys(imports[path]).forEach(function (importedSymbol) {
        rule.push(_postcss2['default'].decl({
          value: importedSymbol,
          prop: imports[path][importedSymbol],
          raws: { before: '\n  ' },
          _autoprefixerDisabled: true
        }));
      });
    });
  };
});

exports['default'] = processor;
module.exports = exports['default'];
/*match*/ /*match*/