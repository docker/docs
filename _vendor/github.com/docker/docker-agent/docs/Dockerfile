FROM jekyll/jekyll:4

COPY Gemfile Gemfile.lock* ./
RUN bundle install

EXPOSE 4000
CMD ["jekyll", "serve", "--host", "0.0.0.0", "--livereload"]
