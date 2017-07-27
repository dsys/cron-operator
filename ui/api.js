const JOBS_LIST_URL = "/k8s-api/apis/batch/v1/jobs";
const CRON_JOBS_LIST_URL = "/k8s-api/apis/kubeheads.pavlov.ai/v1/cronjobs";

export async function fetchJobs() {
  const res = await fetch(JOBS_LIST_URL);
  const resBody = await res.json();
  return resBody.items.reverse();
}

export async function fetchCronJobs() {
  const res = await fetch(CRON_JOBS_LIST_URL);
  const resBody = await res.json();
  return resBody.items.reverse();
}
