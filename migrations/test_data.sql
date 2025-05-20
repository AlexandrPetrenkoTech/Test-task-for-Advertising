INSERT INTO adverts (name, description, price, created_at) VALUES
                                                               ('iPhone 14',     'Brand new iPhone 14, unopened box.',             999.00, NOW() - INTERVAL '5 days'),
                                                               ('Gaming Laptop', 'High-performance laptop for gaming and work.', 1299.00, NOW() - INTERVAL '3 days'),
                                                               ('Coffee Table',  'Wooden coffee table in good condition.',          49.00, NOW() - INTERVAL '7 days'),
                                                               ('Mountain Bike', 'Lightweight bike, perfect for off-road trails.',  300.00, NOW() - INTERVAL '1 day'),
                                                               ('Desk Lamp',     'Minimalist lamp, LED, adjustable arm.',           20.00, NOW());

INSERT INTO photos (advert_id, url, position) VALUES
                                                  (1, 'https://example.com/iphone.jpg', 1),
                                                  (2, 'https://example.com/laptop.jpg', 1),
                                                  (3, 'https://example.com/table.jpg', 1),
                                                  (4, 'https://example.com/bike1.jpg', 1),
                                                  (4, 'https://example.com/bike2.jpg', 2),
                                                  (5, 'https://example.com/lamp.jpg', 1);