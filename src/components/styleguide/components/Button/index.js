import React, { Component } from 'react';
import { DockerFlatIcon, Button } from 'common';
import asExample from '../../asExample';

import mdHeader from './header.md';
import mdApi from './api.md';

import { variants } from 'lib/constants';
import { map } from 'lodash';

@asExample(mdHeader, mdApi)
export default class ButtonDoc extends Component {
  render() {
    return (
      <div>
        <h3>Variations</h3>
        <div>
          { map(variants, (variant) =>
            <Button
              variant={variant}
              key={`${variant}1`}
            >{variant}</Button>
          )}
        </div>
        <div>
          { map(variants, (variant) =>
            <Button
              variant={variant}
              key={`${variant}2`}
              outlined
            >{variant}</Button>
          )}
        </div>
        <div>
          { map(variants, (variant) =>
            <Button
              variant={variant}
              key={`${variant}3`}
              text
            >Link</Button>
          )}
        </div>
        <div>
          { map(variants, (variant) =>
            <Button
              variant={variant}
              key={`${variant}4`}
              disabled
            >{variant}</Button>
          )}
        </div>
        <div>
          { map(variants, (variant) =>
            <Button
              variant={variant}
              key={`${variant}5`}
              outlined
              disabled
            >{variant}</Button>
          )}
        </div>
        <div>
          { map(variants, (variant) =>
            <Button
              variant={variant}
              key={`${variant}6`}
              text
              disabled
            >{variant}</Button>
          )}
        </div>
        <h3>Icon Button</h3>
        <div>
          { map(variants, (variant) =>
            <Button
              variant={variant}
              key={`${variant}7`}
              icon
            >
              <DockerFlatIcon />
            </Button>
          )}
        </div>
        <div>
          { map(variants, (variant) =>
            <Button
              variant={variant}
              key={`${variant}8`}
              icon
              disabled
            >
              <DockerFlatIcon />
            </Button>
          )}
        </div>
        <h3>With Icons (left)</h3>
        <div>
          <div>
            { map(variants, (variant) =>
              <Button
                variant={variant}
                key={`${variant}9`}
                icon="left"
              >
                <DockerFlatIcon />
                Dockerhub
              </Button>
            )}
          </div>
          <div>
            { map(variants, (variant) =>
              <Button
                variant={variant}
                key={`${variant}10`}
                icon="left"
                outlined
              >
                <DockerFlatIcon />
                Dockerhub
              </Button>
            )}
          </div>
        </div>
        <h3>With Icons (right)</h3>
        <div>
          <div>
            { map(variants, (variant) => <Button
              variant={variant}
              key={`${variant}11`}
              icon="right"
            >
              <DockerFlatIcon />
              Dockerhub
            </Button>
            )}
          </div>
          <div>
            { map(variants, (variant) => <Button
              variant={variant}
              key={`${variant}12`}
              icon="right"
              outlined
            >
              <DockerFlatIcon />
              Dockerhub
            </Button>
            )}
          </div>
        </div>
        <h3>As link</h3>
        <div>
          <div>
            { map(variants, (variant) => <Button
              variant={variant}
              element="a"
              key={`${variant}13`}
              icon="right"
              href="//hub.docker.com"
              target="_blank"
            >
              <DockerFlatIcon />
              Dockerhub
            </Button>
            )}
          </div>
          <div>
            { map(variants, (variant) => <Button
              variant={variant}
              element={<a />}
              key={`${variant}13`}
              icon="right"
              href="//hub.docker.com"
              outlined
              target="_blank"
            >
              <DockerFlatIcon />
              Dockerhub
            </Button>
            )}
          </div>
          <div>
            { map(variants, (variant) => <Button
              variant={variant}
              style={{ width: '130px' }}
              element={<a />}
              key={`${variant}13`}
              href="//hub.docker.com"
              text
              target="_blank"
            >
              Dockerhub
            </Button>
            )}
          </div>
        </div>

        <h3>Inverted</h3>
        <div style={{ minHeight: '50px', backgroundColor: '#1AAAF8' }}>
          <Button inverted>Inverted</Button>
          <Button inverted icon>
            <DockerFlatIcon />
          </Button>
          <Button inverted text>Inverted</Button>
          <Button inverted text icon="left">
            <DockerFlatIcon />
            Inverted
          </Button>
          <Button inverted text icon="right">
            <DockerFlatIcon />
            Inverted
          </Button>
          <Button inverted outlined icon="right">
            <DockerFlatIcon />
            Inverted
          </Button>
          <Button inverted disabled>Disabled</Button>
          <Button inverted disabled outlined>Disabled</Button>
        </div>
      </div>
    );
  }
}
