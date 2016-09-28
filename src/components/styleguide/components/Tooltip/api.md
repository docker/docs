
```
import { Tooltip } from 'common';

const tooltipContent = <div>Hi I am a tooltip!!</div>;
const tooltipWithIcon = <DockerCloudIcon size="large" variant="secondary" />;
<Tooltip
  content={tooltipContent}
  placement="top"
>
  <div>
    My Tooltip is on Top
  </div>
</Tooltip>

<Tooltip
  content={tooltipWithIcon}
  placement="bottom"
  trigger={['click']}
>
  <div>
    Click me!
  </div>
</Tooltip>

<Tooltip
  content={tooltipWithIcon}
  placement="top"
  trigger={['focus']}
>
  <input placeholder="Focus on me!" />
</Tooltip>

<Tooltip
  content={tooltipContent}
  placement="top"
  theme="dark"
>
  <div>
    My Tooltip has a Dark Theme!
  </div>
</Tooltip>
```

### Props

| name     | propType                         | default | required | description                                                 |
|----------|----------------------------------|---------|----------|-------------------------------------------------------------|
| align      | Object: alignConfig of  [dom-align](https://github.com/yiminghe/dom-align)                          |        |    no    |  value will be merged into placement's align config. note: can only accept offset and targetOffset                                              
| className      | string                           |        |    no    |  additional className added to tooltip                                                 |
| children      | node                           |        |    yes    |  Item to render the tooltip around                                                    |
| content      | node                           |     |    yes    |  tooltip content - can be any react element or node                                                           |
| mouseEnterDelay     | number                           |   0.1     |    no    |  delay time to show tooltip when mouse enters (sec)                                                 |
| mouseLeaveDelay     | number                           |   0.1     |    no    |  delay time to show tooltip when mouse leaves (sec)                                                 |
| placement     | string                           |   "top"    |    no    |  Where the tooltip should render relative to the children: one of ['left', 'right', 'top', 'bottom', 'topLeft', 'topRight', 'bottomLeft', 'bottomRight']                                        |
| theme       | string         | "white" | no | Possible style / theme for the tooltip. One of "dark" or "white"
| trigger       | string[]         | ['hover'] | no | Which actions cause the tooltip to show: available values are 'hover','click','focus'

---
