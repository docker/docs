{% capture green-check %}![yes](/engine/installation/images/green-check.svg){: style="height: 14px; display: inline-block"}{% endcapture %}
{% capture superscript-link %}[1](#edge-footnote){: style="vertical-align: super; font-size: smaller;" }{% endcapture %}
{: style="width: 75%" }

| Month     | Docker CE Edge                          | Docker CE Stable  | 
|:----------|:----------------------------------------|:------------------|
| January   | {{ green-check }}                       |                   |
| February  | {{ green-check }}                       |                   |
| March     | {{ green-check }}{{ superscript-link }} | {{ green-check }} |
| April     | {{ green-check }}                       |                   |
| May       | {{ green-check }}                       |                   |
| June      | {{ green-check }}{{ superscript-link }} | {{ green-check }} |
| July      | {{ green-check }}                       |                   |
| August    | {{ green-check }}                       |                   |
| September | {{ green-check }}{{ superscript-link }} | {{ green-check }} |
| October   | {{ green-check }}                       |                   |
| November  | {{ green-check }}                       |                   |
| December  | {{ green-check }}{{ superscript-link }} | {{ green-check }} |

`1`: On Linux distributions, these releases will only appear in the `stable`
     channels, not the `edge` channels. For that reason, on Linux distributions,
     you need to enable both channels.
{: id="edge-footnote" }
