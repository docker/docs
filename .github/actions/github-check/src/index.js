const { Octokit } = require("@octokit/rest");
const program = require("commander");

const run = (token, message, details_url) => {
  const [owner, repo] = process.env.GITHUB_REPOSITORY.split("/");
  const head_sha = process.env.GITHUB_SHA;

  const octokit = new Octokit(token);

  return octokit.checks
    .create({
      owner,
      repo,
      name: message,
      head_sha,
      conclusion: "success",
      details_url,
      status: "completed",
      output: {
        title: message
      }
    })
    .catch(err => {
      console.error(err.message);
      process.exit(1);
    });
};

program.arguments("<token> <message> <url>").action(run);

program.parse(process.argv);
