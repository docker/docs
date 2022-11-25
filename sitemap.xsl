<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="2.0"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                xmlns:sitemap="http://www.sitemaps.org/schemas/sitemap/0.9"
                xmlns:image="http://www.google.com/schemas/sitemap-image/1.1"
                xmlns:html="http://www.w3.org/TR/REC-html40">
  <xsl:output version="1.0" method="html" encoding="UTF-8" />
  <xsl:template match="/">
    <html xmlns="http://www.w3.org/1999/xhtml">
      <head>
        <title>Docker Docs XML Sitemap</title>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <style type="text/css">
          html, body {
            font-family: Arial, sans-serif;
            font-size: 16px;
            line-height: 1.0;
          }
          a {
            color: #000;
          }
          .intro {
            background-color: #CFEBF7;
            border: 1px #2580B2 solid;
            padding: 5px 13px 5px 13px;
            margin: 10px;
            width: 800px;
          }
          .intro p {
            line-height: 0.5;
          }
          .list td, .list th {
            padding: 5px;
            font-size: 13px;
            line-height: 1.5;
          }
          .list th {
            text-align: left;
          }
          tr.high {
            background-color: whitesmoke;
          }
        </style>
      </head>
      <body>
        <h1>Docker Docs XML Sitemap</h1>
        <div class="intro">
          <p>This is an XML Sitemap which is supposed to be processed by search engines like <a href="https://www.google.com">Google</a>, <a href="https://www.bing.com">Bing</a>, ...</p>
          <p>You can find more information about XML sitemaps on <a href="https://sitemaps.org/">sitemaps.org</a> and Google's <a href="https://code.google.com/archive/p/sitemap-generators/wikis/SitemapGenerators.wiki">list of sitemap programs</a>.</p>
          <p>This sitemap contains <strong><xsl:value-of select="count(sitemap:urlset/sitemap:url)"/></strong> URLs.</p>
        </div>
        <div id="content">
          <table class="list" cellpadding="5">
            <tr style="border-bottom:1px black solid;">
              <th>Location</th>
              <th>Last Modification</th>
            </tr>
            <xsl:for-each select="sitemap:urlset/sitemap:url">
              <tr>
                <xsl:if test="position() mod 2 != 1">
                  <xsl:attribute  name="class">high</xsl:attribute>
                </xsl:if>
                <td>
                  <xsl:variable name="itemURL">
                    <xsl:value-of select="sitemap:loc"/>
                  </xsl:variable>
                  <a href="{$itemURL}">
                    <xsl:value-of select="sitemap:loc"/>
                  </a>
                </td>
                <td>
                  <xsl:value-of select="concat(substring(sitemap:lastmod,0,11),concat(' ', substring(sitemap:lastmod,12,5)))"/>
                </td>
              </tr>
            </xsl:for-each>
          </table>
        </div>
      </body>
    </html>
  </xsl:template>
</xsl:stylesheet>
