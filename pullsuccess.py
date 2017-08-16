import urllib2
from bs4 import BeautifulSoup
from subprocess import call

call(["rm", "-rf", "kb"])
call(["mkdir", "kb"])
requrl = "https://success.docker.com/@api/deki/pages/"
response = urllib2.urlopen(requrl)
soup = BeautifulSoup(response, 'lxml')
for page in soup.find_all('page'):
    fileout = '';
    skipme = False
    if str(type(page.path.string)) == "<class 'bs4.element.NavigableString'>":
        if page.path.string.find("Internal") > -1:
            skipme = True
            print 'Skipping ' + page['id'] + ' because path=' + page.path.string
    if skipme == False:
        requrl = 'https://success.docker.com/@api/deki/pages/' + page['id']
        pagemetadata = urllib2.urlopen(requrl)
        metadatasoup = BeautifulSoup(pagemetadata, 'lxml')
        fileout += '---\n'
        fileout += 'title: \"' + page.title.string.replace('"','') + '\"\n'
        fileout += 'id: ' + page['id'] + '\n'
        fileout += 'draftstate: ' + page['draft.state'] + '\n'
        fileout += 'deleted: '  + page['deleted'] + '\n'
        fileout += 'source: https://success.docker.com/@api/deki/pages/' + page['id'] + '/contents' + '\n'
        fileout += 'tags:' + '\n'
        for thistag in metadatasoup.find_all('tag'):
            fileout += '- tag: \"' + thistag['value'] + '\"' + '\n'
        fileout += '---' + '\n'
        fileout += '{% raw %}\n'

        requrl = 'https://success.docker.com/@api/deki/pages/' + page['id'] + '/contents'
        pagecontents = urllib2.urlopen(requrl)
        contentsoup = BeautifulSoup(pagecontents, 'html.parser')
        rawhtml = BeautifulSoup(contentsoup.get_text(), 'html.parser')
        '''
        images = rawhtml.find_all('img')
        for image in images:
            print image
        '''
        fileout += rawhtml.prettify() + '\n'
        fileout += '{% endraw %}\n'
        f = open('kb/' + page['id'] + '.md', 'w+')
        f.write(fileout.encode('utf8'))
        f.close
        print 'Success writing kb/' + page['id'] + '.md'
