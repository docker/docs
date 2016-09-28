module.exports = function indexPageProcessor(){
    return{
        $runBefore:['rendering-docs'],
        $runAfter:['componentsDataProcessor'],
        $process:function(docs){
            docs.push({
                template:'index.template.html',
                outputPath:'index.html',
                path:'index.html'
            })
        }
    }
};
