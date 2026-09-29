-- Correct only the original broad-catalog CTA; leave author customizations intact.
UPDATE homepage_hero SET primary_cta_text = 'View All Products'
WHERE primary_cta_text = 'View Displays' AND primary_cta_url = '/products';
