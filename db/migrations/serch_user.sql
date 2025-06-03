SELECT name, embedding <=> $1::vector as distance, phone, profession, tg
FROM person 
WHERE embedding IS NOT NULL
ORDER BY distance
LIMIT 5;