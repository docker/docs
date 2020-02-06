const { Octokit } = require("@octokit/rest");

const run = () => {
  const [owner, repo] = process.env.GITHUB_REPOSITORY.split("/");
  const octokit = new Octokit({ auth: process.env.GITHUB_TOKEN });

  return octokit.checks.create({
    owner,
    repo,
    name: "Check Created by API",
    head_sha: process.env.GITHUB_SHA,
    status: "completed",
    conclusion: "success",
    details_url: process.env.SITE_URL
  });
};

run()
  .then(() => console.log("done"))
  .catch(err => {
    console.error(err);
    process.exit(1);
  });
