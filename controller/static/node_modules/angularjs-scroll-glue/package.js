Package.describe({
    name: "luegg:angularjs-scroll-glue",
    summary: "Scrolls to the bottom of an element on changes",
    description: "An AngularJs directive that automatically scrolls to the bottom of an element on changes in it's scope.",
    version: "2.0.7",
    git: "https://github.com/Luegg/angularjs-scroll-glue.git"
});

Package.onUse(function(api) {
    api.versionsFrom('METEOR@1.0');

    api.addFiles('src/scrollglue.js', 'client');
});
