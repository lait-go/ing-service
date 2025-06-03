import sys
from sentence_transformers import SentenceTransformer

model = SentenceTransformer('paraphrase-multilingual-MiniLM-L12-v2')

text = sys.argv[1]
embedding = model.encode(text).tolist()

print(",".join(str(x) for x in embedding))
