from fastapi import FastAPI
from  pydantic import BaseModel
from fastapi.middleware.cors import CORSMiddleware
from app.payment import Payment

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],  # your frontend URL
    allow_credentials=True,
    allow_methods=["*"],  # allows all HTTP methods
    allow_headers=["*"],  # allows all headers
)



@app.get('/')
def main():
    return 'hi'


class Payit(BaseModel):
    email: str = 'test@gmail.com'
    fname: str = 'test'
    lname: str = 'lastname'
    amount: int
    rdurl: str


@app.post('/pay')
def pay(payit: Payit):
    data = Payment.pay(email=payit.email, fname=payit.fname,
                       lname=payit.lname, amount=payit.amount, rdurl=payit.rdurl)
    return data


class Txnum(BaseModel):
    ref_num: str = ''


@app.post('/verify')
def verify(txnum: Txnum):
    ver = Payment.verify(tx_num=txnum.ref_num)
    return ver