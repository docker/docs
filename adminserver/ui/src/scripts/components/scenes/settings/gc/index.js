'use strict';

import React, { Component, PropTypes } from 'react';
import ButtonSelect from 'components/common/buttonSelect';
import Button from 'components/common/button';
import Dropdown from 'components/common/dropdown';
import Input from 'components/common/input';
import classNames from 'classnames';
import buttonSelectStyles from 'components/common/buttonSelect/buttonSelect.css';
import FontAwesome from 'components/common/fontAwesome';
import Spinner from 'components/common/spinner';
import autoaction from 'autoaction';
import ui from 'redux-ui';
import moment from 'moment-timezone';
import css from 'react-css-modules';
import {
    getLastGCSavings,
    getGCSchedule,
    updateGCSchedule,
    deleteGCSchedule,
    runGC,
    getGCStatus,
    stopGC
} from 'actions/settings';
import { addBannerNotification, addGrowlNotification, removeBannerNotification } from 'actions/notifications';
import consts from 'consts';
import { createValidator, required, onlyIf } from 'validation';
import { createStructuredSelector } from 'reselect';
import { gcSchedule, gcTimeout, gcLastRun, gcRunning } from 'selectors/settings';
import { connect } from 'react-redux';
import { mapActions } from 'utils';
import { reduxForm } from 'redux-form';
import styles from './gc.css';

const mapState = createStructuredSelector({
    schedule: gcSchedule,
    timeout: gcTimeout,
    lastRun: gcLastRun,
    running: gcRunning
});

const DO_NOT_REPEAT = 'DO_NOT_REPEAT';

// This component loads data for the GC schedule and renders the UI once ready
@autoaction({
    getGCSchedule: [],
    getLastGCSavings: []
}, {
    getGCSchedule,
    getLastGCSavings
})
export default class GCLoader extends Component {
    render() {
        const status = [
            [consts.settings.GET_GC_SCHEDULE]
        ];

        return (
            <Spinner loadingStatus={ status }>
                <GarbageCollection { ...this.props } />
            </Spinner>
        );
    }
}



@connect(
    mapState,
    mapActions({
        updateGCSchedule,
        deleteGCSchedule,
        runGC,
        getGCStatus,
        stopGC,
        getLastGCSavings,
        addBannerNotification,
        addGrowlNotification,
        removeBannerNotification
    })
)
@ui({
    state: {
        startAfterSubmit: false,
        selected: (props) => {
            if (props.timeout !== '0') {
                return consts.settings.TIMEOUT;
            } else if (props.schedule) {
                return consts.settings.UNTIL_DONE;
            } else {
                return consts.settings.NEVER;
            }
        },
        pollingIntervalID: 0,
        pollInProgress: false
    }
})
@reduxForm({
    form: 'gcschedule',
    fields: [
      'timeout',
      'schedule',

      // this field is only used for validation, because redux-form does not pass state/props
      'neverTab'
    ],
    validate: createValidator({
        'schedule': [onlyIf(data => !data.neverTab, required)]
    })
}, (state, props) => {
    // props.timeout is in seconds, but we show the timeout in minutes
    let initialTimeout = props.timeout;
    let initialSchedule = props.schedule;

    // set a sample timeout of 1 minute to populate the input with, even if the 'Timeout' option is not selected
    if (!initialTimeout || initialTimeout === '0') {
        initialTimeout = '1';
    } else {
        let parsed = parseInt(initialTimeout);
        initialTimeout = '' + Math.round(parsed / 60);

        if (!initialSchedule) {
            initialSchedule = DO_NOT_REPEAT;
        }
    }

    return {
      initialValues: {
        timeout: initialTimeout,
        schedule: initialSchedule,
        neverTab: props.ui.selected === consts.settings.NEVER
      }
    };
})
// this action will trigger componentWillReceiveProps and start polling GC if necessary on load
@autoaction({
    getGCStatus: []
}, {
    getGCStatus
})
@css(styles)
export class GarbageCollection extends Component {

    static propTypes = {
        ui: PropTypes.object,
        updateUI: PropTypes.func,
        schedule: PropTypes.string,
        timeout: PropTypes.string,
        lastRun: PropTypes.object,
        running: PropTypes.bool,
        actions: PropTypes.object,
        fields: PropTypes.object,
        handleSubmit: PropTypes.func
    };

    componentWillReceiveProps(nextProps) {
        if (!this.props.running && nextProps.running) {
            // if we went from not running to running, start polling
            this.startGCPolling();
        }
    }

    chooseScene = (sceneChoice) => {
        const {
            fields: { neverTab }
        } = this.props;

        neverTab.value = sceneChoice === consts.settings.NEVER;
        this.props.updateUI({ selected: sceneChoice });
    }

    startGCPolling = () => {
        const {
            actions,
            updateUI,
            ui: { pollInProgress }
        } = this.props;

        if (pollInProgress) {
            return;
        }

        let doneGCPolling = this.doneGCPolling;

        updateUI({
            pollingIntervalID: setInterval(function() {
                actions.getGCStatus().then(function(resp) {
                    if (!(resp.data.registryGC.status === 'running')) {
                        doneGCPolling();
                    }
                });
            }, consts.settings.GC_POLLING_INTERVAL),
            pollInProgress: true
        });
    }

    doneGCPolling = () => {
        const {
            actions,
            ui: { pollingIntervalID },
            updateUI
        } = this.props;

        clearInterval(pollingIntervalID);

        updateUI({
            pollInProgress: false
        });

        // refresh the lastSavings, since it has probably updated
        actions.getLastGCSavings();

        // notify
        actions.removeBannerNotification('GC_IN_PROGRESS');
        actions.addGrowlNotification({
            title: 'Garbage collection done!',
            message: 'Can push images again',
            status: 'success',
            img: '/public/img/broom-white.png'
        });
    }

    onSubmit = (data) => {
        const {
            actions,
            ui: { selected, startAfterSubmit },
            updateUI
        } = this.props;

        const {
            NEVER,
            TIMEOUT
        } = consts.settings;

        if (selected !== TIMEOUT) {
            // run infinitely if there is no timeout set
            data.timeout = 0;
        }

        // if do not repeat is selected, we should still update the timeout
        // also, if never run is selected then erase the schedule
        if (data.schedule === DO_NOT_REPEAT || selected === NEVER) {
            data.schedule = '';
        }

        let parsedTimeoutMinutes = parseInt(data.timeout);

        // api expects timeout in seconds
        data.timeout = parsedTimeoutMinutes * 60;

        actions.updateGCSchedule(data);

        if (startAfterSubmit) {
            updateUI({ startAfterSubmit: false });
            actions.runGC();
            actions.addBannerNotification({
                id: 'GC_IN_PROGRESS',
                class: 'warning',
                message: 'Deleting images... No one can push images while we clean your storage. Go to Garbage Collection.',
                img: '/public/img/broom-white.png',
                url: '/admin/settings/gc'
            });
            this.startGCPolling();
        }
    }

    // called right before form onSubmit when 'Save & start' is clicked
    startGC = () => {
        this.props.updateUI({ startAfterSubmit: true });
    }

    // takes an ISO8601 formatted timestring and returns time and date strings
    // eg. 2016-05-23T18:02:07.183539695Z becomes:
    // timeString: 11:02am PDT
    // dateString: Mon, 23 May 2016
    getTimeAndDateWithTimezone = (time) => {
        let parsedDateTime = moment(time);
        let timeZone = moment.tz.guess();
        let tzParsedDateTime = parsedDateTime.tz(timeZone);
        let timeString = tzParsedDateTime.format('h:mma z');
        let dateString = parsedDateTime.format('ddd, D MMMM YYYY');
        return {
            timeString: timeString,
            dateString: dateString
        };
    }

    render() {

        const {
            actions,
            fields: { schedule, timeout },
            handleSubmit,
            lastRun,
            running,
            ui: { selected }
        } = this.props;

        const {
            NEVER,
            TIMEOUT,
            UNTIL_DONE
        } = consts.settings;

        let lastRunTimes = {};
        if (lastRun.layers) {
            lastRunTimes = this.getTimeAndDateWithTimezone(lastRun.time);
        }

        return (
            <div>
                <div styleName='banner'>
                    <img styleName='bannerImg' src='/public/img/broom.png' alt='' />
                    <div styleName='bannerText'>
                        <h2 styleName='bannerHeader'>Remove Untagged Images</h2>
                        <p styleName='bannerBody'>
                            Run garbage collection on your storage backend to remove deleted tags and images. <a href='//docs.docker.com/docker-trusted-registry/repos-and-images/delete-images/'>Learn more <FontAwesome icon='fa-external-link' /></a>
                        </p>
                    </div>
                </div>

                { running ?
                    (
                        <div styleName='collectionInProgress'>
                            <h3 styleName='sectionTitle'>Deleting images...</h3>
                            <div>
                                <FontAwesome styleName='spinning' animate={ 'spin' } icon={ 'fa-repeat' } /><span styleName='deletingText'>Sit back and relax while we clean up your storage.</span>
                            </div>
                            <Button variant='alert' onClick={ actions.stopGC }>Stop</Button>
                        </div>
                    ) : (
                        <div>
                            <div>
                                <h3 styleName='sectionTitle'>Delete images</h3>
                            </div>
                            <form method='POST' onSubmit={ handleSubmit(::this.onSubmit) }>
                                <ButtonSelect
                                    initialChoice={ selected }
                                    onChange={ ::this.chooseScene }>
                                    <div icon='fa-hourglass-end' primaryText='Until done' secondaryText='This may take a while' value={ UNTIL_DONE }></div>
                                    <div value={ TIMEOUT } styleName='gcTimeoutButton'>
                                        <div styleName='icon'><FontAwesome icon={ 'fa-hourglass-half' } /></div>
                                        <div styleName='form'>
                                            <div styleName='primaryText'>
                                                For
                                            </div>
                                            <div styleName='inputContainer'>
                                                <Input type={ 'number' } placeholder={ '1' } min={ 1 } formfield={ timeout } markup={ false } />
                                            </div>
                                            <div className={ classNames(buttonSelectStyles.primaryText, styles.primaryText) }>
                                                minute
                                            </div>
                                        </div>
                                    </div>
                                    <div icon='fa-hourglass-o' primaryText='Never' secondaryText='Disable garbage collection' value={ NEVER }></div>
                                </ButtonSelect>
                                { selected !== NEVER && (
                                    <GCScheduler schedule={ schedule } />
                                )
                                }
                                <div styleName='saveButtons'>
                                    { selected !== NEVER && (
                                        <Button onClick={ ::this.startGC }>
                                            <div styleName='buttonText'>Save & Start</div>
                                        </Button>
                                    )
                                    }
                                    <Button variant={ selected === NEVER ? 'primary' : 'primary outline' } type='submit'>
                                        <div styleName='buttonText'>Save</div>
                                    </Button>
                                </div>
                                { // we seem to get '0001-01-01T00:00:00Z' as the lastRun time if GC has not run...
                                    lastRun.layers && lastRun.time !== '0001-01-01T00:00:00Z' && (
                                    <div>
                                        <hr />
                                        <div styleName='layersDeletedSection'>
                                            <div styleName='lastRanSection'>
                                                <div styleName='icon'><FontAwesome icon={ 'fa-clock-o' } /></div>
                                                <div styleName='textSection'>
                                                    <h3 className={ classNames(styles.lastRan, styles.sectionTitle) }>Last ran</h3>
                                                    <div styleName='time'>{ lastRunTimes.timeString }</div>
                                                    <div styleName='date'>{ lastRunTimes.dateString }</div>
                                                </div>
                                            </div>
                                            <div styleName='layersDeleted'>
                                                <h3 className={ classNames(styles.layersDeletedTitle, styles.sectionTitle) }>Layers deleted</h3>
                                                { Object.keys(lastRun.layers).map(function(sha) {
                                                    return <div key={ sha } styleName='layer'><span styleName='sha'>{ sha }</span></div>;
                                                }, this) }
                                            </div>
                                        </div>
                                    </div>
                                )
                                }
                            </form>
                        </div>
                    ) }
            </div>
        );
    }
}

@css(styles)
class GCScheduler extends Component {

    static propTypes = {
        schedule: PropTypes.object
    }

    render() {

        const schedule = this.props.schedule;

        return (
            <div>
                <h3 styleName='sectionTitle'>Repeat</h3>
                <div className={ classNames(styles.row, styles.scheduleForm) }>
                    <div styleName='gcScheduleSelect'>
                        <Dropdown values={ schedule }>
                            <option value=''>Custom cron schedule</option>
                            <option selected={ schedule.value === '0 * * *' } value='0 * * *'>Daily at midnight</option>
                            <option selected={ schedule.value === '1 * * 6' } value='1 * * 6'>Every Saturday at 1AM</option>
                            <option selected={ schedule.value === '1 * * 0' } value='1 * * 0'>Every Sunday at 1AM</option>
                            <option selected={ schedule.value === DO_NOT_REPEAT } value={ DO_NOT_REPEAT }>Do not repeat</option>
                        </Dropdown>
                    </div>
                    {
                        schedule.value !== DO_NOT_REPEAT && (
                            <div styleName='gcScheduleInput'>
                                <Input type='text'
                                placeholder='Hour, Day of Mo, Mo, Weekday'
                                formfield={ schedule }/>
                            </div>
                        )
                    }
                </div>
            </div>
        );
    }
}
