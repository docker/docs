describe('ngTableFilterConfig', function () {
    var ngTableFilterConfig,
        ngTableFilterConfigProvider;

    beforeEach(module("ngTable"));
    beforeEach(function(){
        module(function(_ngTableFilterConfigProvider_){
            ngTableFilterConfigProvider = _ngTableFilterConfigProvider_;
            ngTableFilterConfigProvider.resetConfigs();
        });
    });

    beforeEach(inject(function (_ngTableFilterConfig_) {
        ngTableFilterConfig = _ngTableFilterConfig_;
    }));

    describe('setConfig', function(){{

        it('should set aliasUrls supplied', function(){
            ngTableFilterConfigProvider.setConfig({
                aliasUrls: {
                    'text': 'custom/url/custom-text.html'
                }
            });

            expect(ngTableFilterConfig.config.aliasUrls.text).toBe('custom/url/custom-text.html');
        });

        it('should merge aliasUrls with previous values', function(){
            ngTableFilterConfigProvider.setConfig({
                aliasUrls: {
                    'text': 'custom/url/text.html'
                }
            });

            ngTableFilterConfigProvider.setConfig({
                aliasUrls: {
                    'number': 'custom/url/custom-number.html'
                }
            });

            expect(ngTableFilterConfig.config.aliasUrls.text).toBe('custom/url/text.html');
            expect(ngTableFilterConfig.config.aliasUrls.number).toBe('custom/url/custom-number.html');
        });
    }});

    describe('getTemplateUrl', function(){

        it('explicit url supplied', function(){
            var explicitUrl = 'path/to/my-template.html';
            expect(ngTableFilterConfig.getTemplateUrl(explicitUrl)).toBe(explicitUrl);
        });

        it('inbuilt alias supplied', function(){
            expect(ngTableFilterConfig.getTemplateUrl('text')).toBe('ng-table/filters/text.html');
        });

        it('custom alias supplied', function(){
            expect(ngTableFilterConfig.getTemplateUrl('my-template')).toBe('ng-table/filters/my-template.html');
        });

        it('alias registered with custom url', function(){
            ngTableFilterConfigProvider.setConfig({ aliasUrls: {
                'my-template': 'custom/url/my-template.html'
            }});
            expect(ngTableFilterConfig.getTemplateUrl('my-template')).toBe('custom/url/my-template.html');
        });

        it('inbuilt alias registered with custom url', function(){
            ngTableFilterConfigProvider.setConfig({ aliasUrls: {
                'text': 'custom/url/custom-text.html'
            }});
            expect(ngTableFilterConfig.getTemplateUrl('text')).toBe('custom/url/custom-text.html');
        });
    });
});
