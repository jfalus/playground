
import numpy as np
from tqdm import tqdm

CARD_COUNT = 40000
GAMES = 300

cards = np.stack([np.vstack((
    np.sort(np.random.choice(15,5,replace=False)),
    np.sort(np.random.choice(15,5,replace=False))+15,
    np.sort(np.random.choice(15,5,replace=False))+30,
    np.sort(np.random.choice(15,5,replace=False))+45,
    np.sort(np.random.choice(15,5,replace=False))+60,
)).T for _ in range(CARD_COUNT)])

marked_all = np.zeros((GAMES,*cards.shape),dtype=bool)

seqs = np.stack([np.random.permutation(75) for _ in range(GAMES)])

wins = np.zeros(GAMES,dtype=int)

for idx in tqdm(range(75)):
    mask = wins==0
    marked_all|=np.equal(cards,seqs[:,idx].reshape(-1,1,1,1),where=mask.reshape(-1,1,1,1))
    new_row_wins = mask&marked_all.all(axis=2).reshape(GAMES,-1).any(axis=1)
    new_col_wins = mask&marked_all.all(axis=3).reshape(GAMES,-1).any(axis=1)
    wins += new_row_wins + (2*new_col_wins)

row_wins = (wins==1).sum()
col_wins = (wins==2).sum()
both_wins = (wins==3).sum()

print(row_wins,col_wins,both_wins)



# row_wins = 0
# col_wins = 0

# for _ in tqdm(range(GAMES),smoothing=0):
#     seq = np.random.permutation(75)
#     for num in seq:
#         marked|=cards==num
#         if marked.all(axis=1).any():
#             row_wins+=1
#             break
#         if marked.all(axis=2).any():
#             col_wins+=1
#             break
#     marked = np.zeros_like(cards)

# print(f'{row_wins=}\n{col_wins=}')
