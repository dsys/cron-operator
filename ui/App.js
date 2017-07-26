import CronJob from './CronJob';
import Job from './Job';
import JobListFilters from './JobListFilters';
import React from 'react';
import { Intent, Spinner } from '@blueprintjs/core';
import { fetchJobs, fetchCronJobs } from './api';

const REFRESH_INTERVAL = 5000;

function filterJob({ status }, { active, succeeded, failed }) {
  return (
    (active === 0 ||
      (active === 1 && status.active > 0) ||
      (active === -1 && !status.active)) &&
    (succeeded === 0 ||
      (succeeded === 1 && status.succeeded > 0) ||
      (succeeded === -1 && !status.succeeded)) &&
    (failed === 0 ||
      (failed === 1 && status.failed > 0) ||
      (failed === -1 && !status.failed))
  );
}

export default class App extends React.Component {
  state = {
    loading: true,
    jobs: [],
    cronJobs: [],
    filters: { active: 0, succeeded: 0, failed: 0 }
  };

  async componentWillMount() {
    const jobs = await fetchJobs();
    const cronJobs = await fetchCronJobs();
    this.setState({ jobs, cronJobs, loading: false });

    this.refreshInterval = setInterval(async () => {
      this.setState({ loading: true });
      const jobs = await fetchJobs();
      const cronJobs = await fetchCronJobs();
      this.setState({ jobs, cronJobs, loading: false });
    }, REFRESH_INTERVAL);
  }

  componentWillUnmount() {
    clearInterval(this.refreshInterval);
    this.refreshInterval = null;
  }

  handleChangeFilters = diff => {
    this.setState({ filters: { ...this.state.filters, ...diff } });
  };

  render() {
    const { loading, jobs, cronJobs } = this.state;

    if (loading && jobs.length === 0 && cronJobs.length === 0) {
      return (
        <div>
          <Spinner intent={Intent.PRIMARY} />
          <style jsx>{`
            div {
              display: flex;
              justify-content: center;
              height: 100%;
              margin: 80px 0;
            }
          `}</style>
        </div>
      );
    } else {
      const filteredJobs = jobs.filter(j => filterJob(j, this.state.filters));

      return (
        <div className="app">
          <h1>
            All Jobs
            <div className="spinner">
              {loading &&
                <Spinner intent={Intent.PRIMARY} className="pt-small" />}
            </div>
          </h1>
          <h2>
            {cronJobs.length} CronJobs
          </h2>
          {cronJobs.map((cronJobData, i) =>
            <CronJob key={i} data={cronJobData} />
          )}
          <h2>
            {jobs.length} Jobs
          </h2>
          <JobListFilters
            active={this.state.filters.active}
            succeeded={this.state.filters.succeeded}
            failed={this.state.filters.failed}
            onChange={this.handleChangeFilters}
          />
          {filteredJobs.map((jobData, i) => <Job key={i} data={jobData} />)}
          <style jsx>{`
            .app {
              margin: 80px auto;
              max-width: 800px;
              padding: 0 20px;
              width: 100%;
            }

            .spinner {
              display: inline-block;
              margin-left: 10px;
            }

            ul {
              margin-bottom: 20px;
            }

            h1 {
              margin-bottom: 40px;
            }

            h2 {
              margin: 40px 0 20px;
            }
          `}</style>
        </div>
      );
    }
  }
}
