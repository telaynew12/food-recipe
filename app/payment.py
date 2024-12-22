from chapa import Chapa
from datetime import datetime

# Initialize Chapa with the API key
key_api = 'CHASECK_TEST-GdLwnDgIq11gkolK98FSWSjgJSMSamrY'
chapa = Chapa(key_api)

class Payment:
    @staticmethod
    def pay(email, fname, lname, amount, rdurl):
        """
        Initiates the payment process with Chapa.
        """
        now = datetime.now()
        tx_un_num = now.strftime("%m%d%Y%H%M%S")
        tx_id = f'tx_{fname}{tx_un_num}'
        
        # Data for the payment request
        data = {
            'email': email,
            'amount': amount,
            'first_name': fname,
            'last_name': lname,
            'tx_ref': tx_id,
            'return_url': rdurl,
            'customization': {
                'title': '2utube',
                'description': 'Payment for your services',
            }
        }

        # Initialize payment and return the response
        response = chapa.initialize(**data)
        return {'detail': response, 'tx_id': tx_id}

    @staticmethod
    def verify(tx_num):
        """
        Verifies the payment with Chapa.
        """
        response = chapa.verify(tx_num)
        return response
