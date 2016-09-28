var convert = require('swagger-converter');
var fs = require('fs');
var dive = require('dive');

function load(filename) {
    return JSON.parse(fs.readFileSync(filename).toString());
}

var resourceListing = load('./gen/docs.json');
var apiDeclarations = [];

function postProcess(doc) {
    for (var path in doc['paths']) {
        var actions = doc['paths'][path];
        for (var method in actions) {
            var action = actions[method];
            // remove 200 for actions that return non-200 codes (swagger 1.2 limitation)
            for (var code in action['responses']) {
                if (code[0] == '2' && code != '200') {
                    delete action['responses']['200'];
                    break;
                }
            }
            var response200 = action['responses']['200'];
            if (response200 && response200['description'] == 'No response was specified') {
                response200['description'] = '';
                action['responses']['200'] = response200;
            }
        }
    }
    return doc;
}

dive('./gen/docs', function(err, path) {
    apiDeclarations.push(load(path));
}, function(err) {
    if (err) {
        console.log(err);
        return
    }
    var swagger2Document = convert(resourceListing, apiDeclarations);
    swagger2Document = postProcess(swagger2Document);
    console.log(JSON.stringify(swagger2Document));
});
