var path = require('path');

module.exports = {
	context: __dirname, 
    entry: './file-to-annotate',
    output: {
        path: __dirname + '/dist',
        filename: 'build.js'
    },
	resolveLoader: {
    	fallback: path.resolve(__dirname, '../../')
	},
    module: {
	    loaders: [
			{test: /\.js$/, loaders: ['loader']},
    	]
	}
}