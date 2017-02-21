module ProductVersion
	class Generator < Jekyll::Generator
		def generate(site)
			# Use `title` and `product_name` to identify versioned pages
			versioned_pages = site.pages.select {
				|page|
					page.data['title'] &&
						page.data['product_name'] &&
						page.data['product_version']
			}

			# Combine product version information
			post_versions = Hash.new {
				|product, title|
					product[title] = Hash.new()
			}

			versioned_pages.each {
				|page|
					post_versions[page.data['product_name']][page.data['title']] ||= []
					post_versions[page.data['product_name']][page.data['title']] << {
						'product_version' => page.data['product_version'],
						'version_url' => page.url
					}
			}

			# Map product information back into relevant pages
			versioned_pages.each {
				|page|
					page.data['version_info'] = post_versions[page.data['product_name']][page.data['title']]
			}
		end
	end
end
