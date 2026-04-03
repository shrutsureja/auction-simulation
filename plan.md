1. we will start the bidders and set them from the config the number of the bidders

2. then setup the auctions will and start them concurrently

3. bidders will calculate the bid amount and simulate the sleep 
    - send the message to the auction with the bid amount
    - no bid will also be sent with the sleep time

4. The auction will reveice all the bids and will wait for the timeout

5. once the timeout is reached the auction will finalize the winner and will wait until all the bidders have completed their work

6. once all the auction are completed we will calculate the total time taken for the whole process and print the results
