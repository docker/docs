var _ = require('lodash');
module.exports = function componentsDataProcessor() {
    return {
        $runBefore: ['rendering-docs'],
        $runAfter: ['paths-computed'],
        $process: function (docs) {

            var apiParts = [
                {type: 'module', items: []},
                {type: 'service', items: []},
                {type: 'directive', items: []},
                {type: 'object', items: []}
            ];

            docs.forEach(function (doc) {
                var part = _.find(apiParts, function (part) {
                    if (part.type === doc.docType)
                        return true;
                });

                if(part) {
                    part.items.push({
                        name:doc.name,
                        url:doc.outputPath
                    });
                }
            });

            docs.filter(function(doc){
                return doc.methods != null;
            }).forEach(function(doc){
                doc.methods = _.sortBy(doc.methods, "name");
            });

            var navigation = [
                {
                    title: 'API Docucmentation',
                    name: 'api',
                    url: '/docs/api',
                    items: apiParts
                },
                // note: runnable examples are not yet implemented (instead we link to external demo site)
                {
                    title: 'Examples',
                    name: 'examples',
                    url: '/docs/examples',
                    items: []
                }
            ];

            docs.push({
                template: 'content.template.html',
                outputPath: 'partials/content.html',
                path: 'partials/content.html'
            });

            docs.push({
                template: 'api-index.template.html',
                outputPath: 'partials/api-index.html',
                path: 'partials/api-index.html'
            });

            docs.push({
                template: 'get-started.template.html',
                outputPath: 'partials/get-started.html',
                path: 'partials/get-started.html'
            });

            docs.push({
                template: 'nav.template.html',
                outputPath: 'partials/nav.html',
                path: 'partials/nav.html'
            });

            docs.push({
                name: 'NAVSERVICE',
                template: 'constants.template.js',
                outputPath: 'app/js/nav-service.js',
                path: 'app/js/nav-service.js',
                items: navigation
            });
        }
    }
};
