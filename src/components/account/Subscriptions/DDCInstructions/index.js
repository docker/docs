import React, { Component, PropTypes } from 'react';
import {
  ClockIcon,
  AmazonWebServicesIcon,
  AzureIcon,
  AppleIcon,
  WindowsTextIcon,
  LinuxIcon,
  BackButtonArea,
  CopyPullCommand,
  Expand,
} from 'common';
import {
  additionalInformation,
  AWS,
  Azure,
  DTRGuide,
  evaluationInstructions,
  Linux,
  prodInstall,
  UCPGuide,
} from './constants.js';
import map from 'lodash/map';
import { DDC_ID } from 'lib/constants/eusa';
import css from './styles.css';
const { func } = PropTypes;

export default class DDCInstructions extends Component {
  static propTypes = {
    showSubscriptionDetail: func.isRequired,
  }

  showSubscriptionDetail = () => {
    this.props.showSubscriptionDetail(DDC_ID);
  }

  // utility method to generate links
  linkTo = ({ url, label }) => {
    return <a href={url} className={css.link} target="_blank">{label}</a>;
  }

  renderEvaluateDDCOn() {
    const {
      toInstallDDC,
      pullCommand,
      scriptTiming,
    } = evaluationInstructions;
    const clock = <ClockIcon className={css.clock} />;
    const toolboxLink = 'https://docs.docker.com/toolbox/overview/';
    const getDockerLink = this.linkTo({
      label: toInstallDDC[1],
      url: toolboxLink,
    });
    const [sysRequirements, scriptDoes, advOptions] = additionalInformation;
    const mkBulletPoints = (text, index) => {
      return <li key={index} className={css.expandContent}>{text}</li>;
    };
    const mkEnvVars = ({ variable, description }, index) => {
      return (
        <li key={index} className={css.expandContentNoBullet}>
          <div className={css.envVar}>{variable}</div>
          <div>{description}</div>
        </li>
      );
    };
    return (
      <div className={css.section}>
        <div>
          <div className={css.sectionTitle}>
            <span className={css.semiBold}>{'Option 1: '}</span>
            Evaluate a single node on your desktop
          </div>
          <div>
            <AppleIcon className={css.logo} />
            <WindowsTextIcon className={css.logo} />
          </div>
          <div className={css.installSteps}>
            <div>1. {getDockerLink}</div>
            <div>2. {toInstallDDC[2]}</div>
            <div>3. {toInstallDDC[3]}</div>
          </div>
          <div className={css.pullCommand}>
            <CopyPullCommand fullCommand={pullCommand} hasInstruction={false} />
          </div>
          <div className={css.licenseHelpText}>{clock} {scriptTiming}</div>
        </div>
        <div>
          <div className={css.listTitle}>Additional Information</div>
          <div className={css.expandWrapper}>
            <Expand title={sysRequirements.title}>
              <ul className={css.list}>
                {map(sysRequirements.bulletPoints, mkBulletPoints)}
              </ul>
            </Expand>
          </div>
          <div className={css.expandWrapper}>
            <Expand title={scriptDoes.title}>
              <ul className={css.list}>
                {map(scriptDoes.bulletPoints, mkBulletPoints)}
              </ul>
            </Expand>
          </div>
          <div className={css.expandWrapper}>
            <Expand title={advOptions.title}>
              <ul className={css.list}>
                <li className={css.expandContentNoBullet}>
                  {advOptions.bulletPoints[0]}
                </li>
                {map(advOptions.environmentVariables, mkEnvVars)}
              </ul>
            </Expand>
          </div>
        </div>
      </div>
    );
  }

  renderGuideLink = ({ url, label }) => {
    return (
      <div key={label} className={css.guideLink}>
        {this.linkTo({ url, label })}
      </div>
    );
  }

  renderDeployOn = () => {
    const guideLinks = [
      UCPGuide,
      DTRGuide,
      prodInstall,
    ];
    return (
      <div className={css.section}>
        <div>
          <div className={css.sectionTitle}>
            <span className={css.semiBold}>{'Option 2: '}</span>
            Deploy on
          </div>
          <div className={css.deployBlocks}>
            <div className={css.deployBlock}>
              <AmazonWebServicesIcon className={css.logo} />
              {this.linkTo({ url: AWS.url, label: AWS.label })}
            </div>
            <div className={css.deployBlock}>
              <AzureIcon className={css.logo} />
              {this.linkTo({ url: Azure.url, label: Azure.label })}
            </div>
            <div className={css.deployBlock}>
              <LinuxIcon className={css.logo} />
              {this.linkTo({ url: Linux.url, label: Linux.label })}
            </div>
          </div>
        </div>
        <div>
          <div className={css.listTitle}>Resources</div>
          {map(guideLinks, this.renderGuideLink)}
        </div>
      </div>
    );
  }

  render() {
    const supportText =
      'Having trouble getting setup? Please contact ';
    const mailTo = 'mailto:support@docker.com?subject=Docker Datacenter';
    const supportLink =
      this.linkTo({ url: mailTo, label: 'support@docker.com' });
    return (
      <div>
        <BackButtonArea
          onClick={this.showSubscriptionDetail}
          text="Docker Datacenter"
        />
        <div className={css.title}>
          Getting Started with Docker Datacenter
        </div>
        {this.renderEvaluateDDCOn()}
        <hr className={css.hr} />
        {this.renderDeployOn()}
        <hr className={css.hr} />
        <div>{supportText} {supportLink}</div>
      </div>
    );
  }
}
