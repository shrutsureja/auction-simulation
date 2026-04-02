# **Candidate Assessment Task: Auction Simulator**

## **Objective**

The goal of this exercise is to design and implement an **Auction Simulator** that runs multiple auctions concurrently, collects bids from simulated bidders, and measures execution times while adhering to resource constraints.

- There are **100 bidders** participating.

- Each auction is generated based on **20 attributes** of an object.

- All bidders receive these 20 attributes and can respond with a bid. It is not necessary that every bidder will provide a response.

- Each auction will run with a **timeout**. Once the timeout is reached, the auction will close, and the winner will be declared from the bids received up to that point.

- A total of **40 auctions** will run **concurrently** (at the same time).

## **Requirements**

- Measure the **time taken** between the start of the first auction and the completion of the last auction.

- Provide a mechanism to **standardize resources** with respect to the **vCPU and RAM available**.

## **Evaluation Criteria**

- **Correctness**: The auction flow works as described (bidders respond, timeouts handled, winner declared).

- **Time Measurement**: Accurate reporting of the start of the first auction and completion of the last auction.

- **Resource Standardization**: A clear approach to standardizing vCPUs and RAM usage.

- **Clarity**: Code and documentation are structured and easy to understand.

**Deliverables:**  
- gitlab repository with working code  
- video showing code running and sample output  
- sample output file for each auction