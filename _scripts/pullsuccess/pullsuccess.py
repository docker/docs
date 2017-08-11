import urllib2
from bs4 import BeautifulSoup

requrl = "https://success.docker.com/@api/deki/pages/"
response = urllib2.urlopen(requrl)
soup = BeautifulSoup(response, 'lxml')
for page in soup.find_all('page'):
    requrl = "https://success.docker.com/@api/deki/pages/" + page['id']
    pagemetadata = urllib2.urlopen(requrl)
    metadatasoup = BeautifulSoup(pagemetadata, 'lxml')
    skipme = False
    if type(page.path.string) == "<class 'bs4.element.NavigableString'>":
        if page.path.string.find("Internal/") > -1:
            skipme = True
    if skipme == False:
        print '---'
        print 'title: \"' + page.title.string.replace('"','') + '\"'
        print 'id: ' + page['id']
        print 'draftstate: ' + page['draft.state']
        print 'deleted: '  + page['deleted']
        print 'tags:'
        for thistag in metadatasoup.find_all('tag'):
            print '- tag: \"' + thistag['value'] + '\"'
        print '---'
