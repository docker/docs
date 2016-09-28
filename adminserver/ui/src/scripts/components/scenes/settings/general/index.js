'use strict';

import React, {
    Component,
    PropTypes
} from 'react';
const {
    object,
    func,
    string
} = PropTypes;
import Button from 'components/common/button';
import { ToggleWithLabel } from 'components/common/toggleSwitch';
import InputLabel from 'components/common/inputLabel';
import Input from 'components/common/input';
import css from 'react-css-modules';
import styles from 'components/scenes/settings/formstyle.css';
import FontAwesome from 'components/common/fontAwesome';
import ui from 'redux-ui';
import autoaction from 'autoaction';
import { removeBannerNotification } from 'actions/notifications';
import {
    getUpdates,
    saveSettings,
    saveLicense,
    autorefreshLicense,
    getLicenseSettings
} from 'actions/settings';
import { reduxForm } from 'redux-form';
import { connect } from 'react-redux';
import {
    getUpdates as getUpdatesSelector,
    getSettings,
    getLicense,
    getSSLCert
} from 'selectors/settings';
import { createStructuredSelector } from 'reselect';
import { mapActions } from 'utils';
import {
    checkRequiredFields,
    reportMissingFields
} from 'validation';
import moment from 'moment';

const mapState = createStructuredSelector({
    updates: getUpdatesSelector,
    settings: getSettings,
    license: getLicense,
    ssl: getSSLCert
});

@autoaction({
    getUpdates: [],
    getLicenseSettings: []
}, {
    getUpdates,
    getLicenseSettings
})
@connect(
    mapState,
    mapActions({
        saveSettings,
        saveLicense,
        autorefreshLicense,
        removeBannerNotification
    })
)
@ui({
    state: {
        showTLS: false,
        license: (props) => props.license
    }
})
@reduxForm({
    form: 'general',
    fields: [
        'dtrHost',
        'authBypassCA',
        'authBypassOU',
        'httpProxy',
        'httpsProxy',
        'noProxy',
        'checkForUpgrades',
  'reportAnalytics',
  'anonymizeAnalytics',
        'releaseChannel',
        'licenseAutoRefresh',
        'webTLSCert',
        'webTLSKey',
        'webTLSCA'
    ],
    validate: (values) => {
        return reportMissingFields(checkRequiredFields(['dtrHost'], values), {});
    }
}, (state, props) => ({
    initialValues: {
        ...props.settings,
        licenseAutoRefresh: props.license.auto_refresh
    }
}))
@css(styles, { allowMultiple: true })
export default class GeneralSettings extends Component {

    static propTypes = {
        ui: object,
        updateUI: func,
        updates: object,
        fields: object,
        handleSubmit: func,
        actions: object,
        license: object,
        ssl: string,

        // all settings are stored here
        settings: object
    }

    // called from the onChange of the input
    // gets the file and then calls validateNewLicense with the contents
    // then saves the license if it's valid
    selectLicense = (evt) => {
        let files = evt.target.files;
        if (files.length < 1) {
            return;
        }
        let reader = new FileReader();
        reader.onload = (ev) => {
            let data = ev.target.result;
            let license = this.validateNewLicense(data);
            if (license) {
                this.saveLicense(license);
            }
        };
        reader.readAsText(files[0]);
    }

    // simple validation to make sure the user selected a valid license file
    validateNewLicense = (data) => {
        let newLicense, isValid;
        try {
            newLicense = JSON.parse(data);
            isValid = 'key_id' in newLicense;
        } catch (_) {
            isValid = false;
        }
        if (!isValid) {
            alert('You have selected an invalid license file');
            return {};
        }
        return newLicense;
    }

    // save a license
    saveLicense = (jsonLicense) => {
        let rmBannerNotification = this.props.actions.removeBannerNotification;
        this.props.actions.saveLicense(jsonLicense).then(function(resp) {
            if (resp.data.is_valid) {
                rmBannerNotification('NO_LICENSE');
            }
        });
    }

    showTLS = (evt) => {
        evt.preventDefault();
        this.props.updateUI({
            showTLS: !this.props.ui.showTLS
        });
    }

    onSubmit = (data) => {
        if (data.checkForUpgrades === true || data.checkForUpgrades === false) {
          // if toggle was swtiched at all
          data.disableUpgrades = !data.checkForUpgrades;
        } else {
          // otherwise, `checkForUpgrades` is 'undefined'
          // and we default to the previous value
          data.disableUpgrades = this.props.settings.disableUpgrades;
        }
        this.props.actions.saveSettings(data);

        // if the license field is changed we should call that action
        if (data.licenseAutoRefresh !== this.props.license.auto_refresh) {
            this.props.actions.autorefreshLicense(data.licenseAutoRefresh);
        }

    }

    render() {

        const {
            fields: {
                dtrHost,
                httpProxy,
                httpsProxy,
                noProxy,
                checkForUpgrades,
                reportAnalytics,
                anonymizeAnalytics,
                licenseAutoRefresh,
                webTLSCert,
                webTLSKey,
                webTLSCA,
                authBypassCA,
                authBypassOU
            },
            license,
            updates,
            handleSubmit,
            settings,
            ui: {
                showTLS
            }
        } = this.props;

        let initialCheckForUpgrades = true;
        if (settings.disableUpgrades) {
          initialCheckForUpgrades = false;
        }

        let upgradeMessage;
        if (updates.upgradeAvailable !== undefined) {
            upgradeMessage = updates.upgradeAvailable ?
              <strong>Your current Trusted Registry version is out of date.<br />
                <a href='https://docs.docker.com/docker-trusted-registry/release-notes/' target='_blank'>Click here to read the release notes</a> for upgrade information.
              </strong> :
              <span>Your Trusted Registry is up-to-date!</span>;
        }

        const licenseExpirationDate = moment(license.expiration).format('MMMM D, YYYY [at] h:mmA');

        return (

            <form method='POST' onSubmit={ handleSubmit(::this.onSubmit) } styleName='generalForm'>

                 <div styleName='row'>
                     <div styleName='formbox halfColumn'>
                         <h2>Updates</h2>
                         <p>{ upgradeMessage }</p>
                         <InputLabel>Version</InputLabel>
                         <p>{ window.currentVersion }</p>
                         <ToggleWithLabel
                             labelOptions={ {
                                 inline: true,
                                 labelText: 'Check for updates',
                                 tip: 'Make outbound connections for update checks. If disabled, you will not be notified when important updates are available.'
                             } }
                             initial={ initialCheckForUpgrades }
                             formField={ checkForUpgrades }
                          /><br />
                         <p><a href='//docs.docker.com/docker-trusted-registry/release-notes/' target='dtrReleaseNotes'>View release notes <FontAwesome icon='fa-external-link' /></a></p>
                     </div>

                     <div styleName='formbox halfColumn'>
                            <h2>License</h2>
                            { license && Object.keys(license).length ?
                                <div>
                                    <InputLabel>Tier</InputLabel>
                                    <p>
                                        { license.tier === 'Production' ?
                                            <FontAwesome icon='fa-trophy' styleName='licenseProd' /> : undefined
                                        }
                                        { license.tier }
                                    </p>
                                    <InputLabel>Expires</InputLabel>
                                    <p>{ licenseExpirationDate }</p>
                                    <div styleName='row'>
                                    <div styleName='halfColumn'>
                                        <InputLabel>License ID</InputLabel>
                                        <p styleName='licenseKey'>{ license.key_id }</p>
                                    </div>
                                    <div styleName='halfColumn'>
                                        <Button
                                          id='apply-license-button'
                                          variant='primary outline'
                                          type='button'
                                          onClick={
                                            () => {
                                              if (navigator.userAgent.toLowerCase().indexOf('firefox') > -1) {
                                                this.refs.newLicense.click();
                                              }
                                            }
                                        }>
                                            <label htmlFor='newLicense' styleName='licenseLabel'>Apply new license</label>
                                            <input ref='newLicense' id='newLicense' type='file' accept='.lic' onChange={ this.selectLicense } styleName='licenseInput' />
                                        </Button>
                                    </div>
                                    </div>
                                    <ToggleWithLabel
                                     initial={ license.auto_refresh }
                                     labelOptions={ {
                                         inline: true,
                                         labelText: 'Auto refresh license'
                                     } }
                                    formField={ licenseAutoRefresh }
                                    />
                                </div>
                                :
                                /* If no license uploaded yet */
                                <div styleName='row'>
                                    <div styleName='halfColumn'>
                                        <InputLabel>You need a license to run Trusted Registry</InputLabel>
                                        <p><a href='//hub.docker.com/enterprise/trial/' target='dtrLicense'>Get one here <FontAwesome icon='fa-external-link' /></a></p>
                                    </div>
                                    <div styleName='halfColumn'>
                                        <Button variant='primary outline' type='button'>
                                            <label htmlFor='newLicense' styleName='licenseLabel'>Apply new license</label>
                                            <input id='newLicense' type='file' accept='.lic' onChange={ this.selectLicense } styleName='licenseInput' />
                                        </Button>
                                    </div>
                                </div>
                            }
                      </div>
                 </div>

                <div styleName='formbox'>
                     <h2>Domain</h2>
                     <div>
                        <InputLabel hint='The public domain name (or IP) and port of the server.'>Load Balancer/Public Address</InputLabel>
                     </div>
                     <div>
                        <Input type='text' placeholder='example.com' formfield={ dtrHost } />
                     </div>
                     <p>
                        <a href='#' onClick={ ::this.showTLS }>
                            { showTLS ?
                                <span><FontAwesome icon='fa-caret-down' /> Hide </span> :
                                <span><FontAwesome icon='fa-caret-right' /> Show </span>
                            }
                            TLS settings
                        </a>
                     </p>
                     { showTLS ?
                       <div styleName='tls'>
                           <InputLabel
                               isOptional={ true }
                               tip={ <span>Certificate issued by a Certificate Authority. If there are any intermediate certificates
                                    they should be included here in the correct order.
                                    You can generate your own certificates for Trusted Registry using a public service or your enterprise&#39;s infrastructure.</span> }>
                                    TLS certificate
                            </InputLabel>
                           <Input isTextarea formfield={ webTLSCert } />
                           <InputLabel
                               isOptional={ true }
                               tip='The key you used to generate your request for a TLS Certificate.'>TLS private key</InputLabel>
                           <Input isTextarea formfield={ webTLSKey } />
                           <InputLabel
                               isOptional={ true }
                               tip='The CA authority used to create your TLS certificate'>TLS CA</InputLabel>
                           <Input isTextarea formfield={ webTLSCA } />
                       </div> :
                       undefined
                     }
                 </div>

                <div styleName='formbox'>
                    <h2>Proxies</h2>
                        <div styleName='row'>
                            <div styleName='halfColumn'>
                                <InputLabel
                                    tip='Proxy server for external HTTP requests. Set by the Docker container.'>HTTP proxy</InputLabel>
                                <p>{ httpProxy.value || 'N/A' }</p>
                            </div>
                            <div styleName='halfColumn'>
                                <InputLabel
                                    tip='Proxy server for external HTTPS requests. Set by the Docker container.'>HTTPS proxy</InputLabel>
                                <p>{ httpsProxy.value || 'N/A' }</p>
                            </div>
                        </div>
                        <InputLabel
                            tip='Proxy bypass for HTTP/S requests.'>Don&apos;t use proxies for...</InputLabel>
                        <p>{ noProxy.value || 'N/A' }</p>
                 </div>

                <div className={ styles.formbox }>
                    <InputLabel
                        tip='The TLS certificate of the Certificate Authority used to verify the client certificate allowed to bypass auth for the registry. This is used primarily by UCP.'
                        isOptional={ true }>Auth Bypass TLS Root CA</InputLabel>
                    <Input isTextarea formfield={ authBypassCA } />

                    <InputLabel
                        isOptional={ true }
                        tip='The OU (organization unit) of the client certificate used to bypass auth for the registry. Leave blank if used with UCP.'
                    >Auth Bypass OU</InputLabel>
                    <Input
                        type='text'
                        placeholder='ucp'
                        // manually pass formfield
                        formfield={ authBypassOU } />
                </div>

                 <div className={ styles.formbox }>
                     <h2>Analytics</h2>
                     <p>Sends periodic usage data reports so Docker can make Docker Trusted Registry better.</p>
                     <ToggleWithLabel
                         labelOptions={ {
                             inline: true,
                             labelText: 'Send data'
                         } }
                         initial={ settings.reportAnalytics }
                         formField={ reportAnalytics }
                      />
      { reportAnalytics ?
        <ToggleWithLabel
          labelOptions={ {
            inline: true,
            labelText: 'Make data anonymous'
          } }
          initial={ settings.anonymizeAnalytics }
          formField={ anonymizeAnalytics }
        />
      : undefined
      }
                 </div>
                 <div styleName='buttonRow'>
                   <Button variant='primary'>Save</Button>
                 </div>
            </form>
        );
    }
}
