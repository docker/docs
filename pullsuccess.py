import urllib2
from bs4 import BeautifulSoup

requrl = "https://success.docker.com/@api/deki/pages/"
response = urllib2.urlopen(requrl)
soup = BeautifulSoup(response, 'lxml')
for page in soup.find_all('page'):
    skipme = False
    if str(type(page.path.string)) == "<class 'bs4.element.NavigableString'>":
        if page.path.string.find("Internal") > -1:
            skipme = True
            print 'Skipping ' + page['id'] + ' because path=' + page.path.string
    if skipme == False:
        requrl = 'https://success.docker.com/@api/deki/pages/' + page['id']
        pagemetadata = urllib2.urlopen(requrl)
        metadatasoup = BeautifulSoup(pagemetadata, 'lxml')
        print '---'
        print 'title: \"' + page.title.string.replace('"','') + '\"'
        print 'id: ' + page['id']
        print 'draftstate: ' + page['draft.state']
        print 'deleted: '  + page['deleted']
        print 'source: https://success.docker.com/@api/deki/pages/' + page['id'] + '/contents'
        print 'tags:'
        for thistag in metadatasoup.find_all('tag'):
            print '- tag: \"' + thistag['value'] + '\"'
        print '---'

        '''
        requrl = 'https://success.docker.com/@api/deki/pages/' + page['id'] + '/contents'
        pagecontents = urllib2.urlopen(requrl)
        contentsoup = BeautifulSoup(pagecontents, 'html.parser')
        rawhtml = BeautifulSoup(contentsoup.get_text(), 'html.parser')
        images = rawhtml.find_all('img')
        for image in images:
            print image
        print rawhtml.prettify()
        '''
