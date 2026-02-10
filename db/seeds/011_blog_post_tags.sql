INSERT INTO blog_post_tags (blog_post_id, blog_tag_id)
SELECT bp.id, bt.id FROM blog_posts bp, blog_tags bt
WHERE bp.slug = 'future-interactive-flat-panels-education' AND bt.slug IN ('ifp', 'education', 'edtech', 'classroom-technology', 'interactive-learning');

INSERT INTO blog_post_tags (blog_post_id, blog_tag_id)
SELECT bp.id, bt.id FROM blog_posts bp, blog_tags bt
WHERE bp.slug = 'new-desktop-series-launch-performance-efficiency' AND bt.slug IN ('desktops', 'enterprise');

INSERT INTO blog_post_tags (blog_post_id, blog_tag_id)
SELECT bp.id, bt.id FROM blog_posts bp, blog_tags bt
WHERE bp.slug = 'step-by-step-guide-setting-up-interactive-display' AND bt.slug IN ('ifp', 'interactive-learning');

INSERT INTO blog_post_tags (blog_post_id, blog_tag_id)
SELECT bp.id, bt.id FROM blog_posts bp, blog_tags bt
WHERE bp.slug = 'bluejay-labs-partners-leading-education-board' AND bt.slug IN ('education', 'edtech');

INSERT INTO blog_post_tags (blog_post_id, blog_tag_id)
SELECT bp.id, bt.id FROM blog_posts bp, blog_tags bt
WHERE bp.slug = 'understanding-ops-module-technology-enterprise' AND bt.slug IN ('ops-modules', 'enterprise');

INSERT INTO blog_post_tags (blog_post_id, blog_tag_id)
SELECT bp.id, bt.id FROM blog_posts bp, blog_tags bt
WHERE bp.slug = 'digital-signage-trends-corporate-communications' AND bt.slug IN ('digital-signage', 'corporate-training');

INSERT INTO blog_post_tags (blog_post_id, blog_tag_id)
SELECT bp.id, bt.id FROM blog_posts bp, blog_tags bt
WHERE bp.slug = 'healthcare-av-solutions-improving-patient-outcomes' AND bt.slug IN ('av-solutions', 'healthcare');
