# Interchain Security


# Introduction

Interchain Security has been referred to in many different ways: Shared Security, Cross Chain Validation, Cross Chain Collateralization, Shared Staking.  This document will restrict use to the following three terms:



*   **Shared Security**
    *   Shared security is a family of technologies that include optimistic rollups, zk-rollups, sharding and interchain security.
*   **Interchain Security**
    *   Interchain Security is the Cosmos specific category of Shared Security that uses IBC (Inter-Blockchain Communication).
*   **Cross Chain Validation**
    *   Cross Chain Validation is the specific IBC level protocol that enables Interchain Security.

While there are many ways that Interchain Security could take place, this document will focus on one instance of Interchain Security that has particularly valuable attributes for the ATOM token and the Cosmos Hub. The resulting technology may be applied to other scenarios with little to no modification, but I will leave those out for now (or only dedicate a small section) since the current priority is to implement this feature for the Cosmos Hub.

At a very high level, Interchain Security is the ability for staking tokens that have been delegated to validators on a Provider Chain to inform the composition of a validator set on a Consumer Chain. Inter-Blockchain Communication is utilized to relay updates of validator stake delegations from the Provider Chain to the Consumer Chain so that the Consumer Chain will have an up-to-date model of which validators can produce blocks on the Consumer Chain. The inclusion of Provider Chain validators can be mandatory or opt-in depending on the requirements of the Consumer Chain. The Provider Chain will honor any proof of validator misbehavior produced by the Consumer Chain as evidence that results in slashing the stake of misbehaving validators on the Provider Chain. In this way the security gained from the value of the stake locked on the Provider Chain will be shared with the Consumer Chain.


# Cosmos Hub User Story

There are two primary reasons that Interchain Security is valuable to the Cosmos Hub. The first reason is because it allows for hub minimalism and the second is to lower the barrier to launching and running secure sovereign decentralized public blockchains. 


## Practical Hub Minimalism

Hub Minimalism is the strategic philosophy that posits that the Cosmos Hub should have as few features as possible in order to decrease the surface area for security vulnerabilities and to reduce the chance of a conflict of interest between user groups. A hub minimalist might be against the governance module from being on the same blockchain as a DEX since users of the governance module must now accommodate users of the DEX even when they have different interests. At best, divergent user groups can peacefully coexist and at worst they may result in hard-forks that diverge in the features of an application.

The current Cosmos Hub is adding more features, which carries some of the risks that hub minimalism is concerned with. Should Interchain Security become available, it would be possible to satisfy  hub minimalists by allowing for each distinct feature of the Cosmos Hub to be an independent chain validated by the same set of ATOM delegated validators. This way the operation of each function could occur independently without affecting the operation of other ATOM secured hub-specific applications.


## Lowering the Barrier to Security

The security of a network is often described as a function of the cost for attacking that network. In Tendermint consensus we target ⅓ and ⅔ of locked stake for various guarantees about liveness and correctness. This means that in order to do any of a variety of attacks against the network, you would need to acquire ⅓+ or ⅔+ of all stake. The crude way to calculate the cost of an attack is to take the quantity of tokens needed to achieve these proportions and multiply it by the current market price for that token. We'll call this the Cost of Corruption.

The Cost of Corruption calculation doesn't account for the availability of any specific token but it does give a very rough estimate for how secure a chain is. It's important that the total value locked (TVL) on a chain remains less than the Cost of Corruption, otherwise the chain should be considered insecure. Since the ability of a chain to serve a valuable purpose is often dependent on the TVL it can handle, it's important to find ways to increase the Cost of Corruption for chains in the Cosmos ecosystem. One method of doing this is to lend the value of the Cosmos ATOM to the Cost of Corruption of any chain. This becomes possible with Interchain Security as the ATOM, which already has a sizeable market cap, can increase the Cost of Corruption of a new chain.


# The Interchain Security Stack

How Interchain Security works at a technical level is still in the process of development but the stack at a high level is well mapped out. It requires new functionality and modifications to current functionality on both Provider Chain and Consumer Chains. The technology can be developed progressively so that a minimum viable set of functionality can be launched as a V1 before an extended set of functionality is launched at a later date.


## Chain Registry

### V1 - Full Validator Set

To work iteratively, the simplest version of Interchain Security is the first milestone and includes the entire validator set of the Provider Chain. In order to ensure that the entire validator set is prepared to validate on a Consumer Chain, it must pass the Provider Chain governance process. The governance proposal is attached to a new Chain Registry module. This module keeps track of Consumer Chains who wish to use the Provider Chain's validator set. Like all governance proposals, the risks and benefits should be socialized off-chain and eventually ratified with an on-chain vote. Should the vote pass, the Consumer Chain will be able to begin the process of using Interchain Security. It is expected that the first Consumer Chains will be thought of as extensions of the Cosmos Hub itself; essentially modules that would otherwise run on the Cosmos Hub directly, but for one reason or another are better suited to be on their own application-specific blockchain. These Consumer Chains can be thought of as having the same security guarantees as the Cosmos Hub, secured by the full force of the ATOM staking token. These Consumer Chains may or may not provide fees and rewards to the Provider Chain validators, depending on the specific application design. They may be deemed valuable enough to the success of Cosmos and the ATOM that validators participate voluntarily, or they may have an application-specific governance token that is used as a fee token and rewarded to validators and their constituent delegators.

## Provider Chain Staking Module

Tendermint uses ABCI to get a set of validators and voting powers from the state machine in order to perform consensus on block production. This information is stored within the staking module of the Cosmos SDK application. In the configuration of Interchain Security, the Consumer Chain also has an instance of Tendermint that uses ABCI to ask for a set of validators and their voting powers. However instead of coming directly from the staking module of the Consumer Chain, in a sense it needs to come from the staking module of the Provider Chain. Practically speaking the state of the staking module of the Provider Chain is relayed periodically via IBC to the Consumer Chain, where it is stored in the Consumer Chain staking module and accessible to Tendermint via ABCI.

### Epoch Staking

The current Staking Module of the Cosmos SDK is moving towards Epoch based staking. This means that instead of validator set delegation amounts being calculated on a per block basis, they will be calculated over some length of time (or blocks) called an Epoch. This will decrease the number of times staking is calculated and generally decrease the complexity involved in staking. The additional complexity of chain relevant stake calculations will similarly benefit from a general simplification of stake calculations, and it will require less packets to be sent between the chains.

## Provider Chain Distribution Module

The distribution module of the Cosmos SDK uses a system called F1 to keep track of how many staking tokens that delegators have bonded to different validators and for how long they've been doing it. Block production rewards and all transaction fees are pooled into the distribution module account at the end of each block, and then distributed to each validator account based on their total voting power. When a delegator wants to withdraw rewards, the distribution module calculates the number of rewards received by their validator on their behalf since the delegator last withdrew, and calculates the outcome based on the amount of stake that belongs to the delegator compared to the total stake of the validator.

Luckily, this system is sufficient to handle the distribution of rewards that come to the Provider Chain from a Consumer Chain. Since the distribution module is collecting all fees from the last epoch, and these fees can be in any variety of denominations, it can be similarly used to distribute Consumer Chain fees. The Consumer Chain can simply use an IBC transfer packet to send the entire batch of fees from a single epoch to the Provider Chain at the end of the epoch, targeted directly at the distribution module account. From the perspective of the distribution module account, these will look like normal fees collected and be distributed to all the validators and their delegators. This simple solution will only work when it is the entire validtor set participating in Interchain Security (V1 Full Validator Set).

## IBC & Cross Chain Validation

There are a number of IBC application layer modules and packets that need to be developed to fully realize the IBC component of Interchain Security. This work has begun with a spec draft from Informal Systems that is visible at [@informalsystems/cross-chain-validation](https://github.com/informalsystems/cross-chain-validation/). Instead of diving into the details of what they are and exactly how they work this section will be reserved for high level responsibilities of these mechanisms.

There are three types of operations within Cross Chain Validation which must be present for Interchain Security to take place:


*   Validator Set Updates
*   Evidence


### Validator Set Updates

The primary duty of Cross Chain Validation is to relay the set of validators and their voting power. The inclusion of a validator within a set relevant to a specific Consumer Chain is designated within the Chain Registry. The voting power denominated in Provider Chain staking token is designated within the staking module. The xStaking module allows for the collation of validators and their delegations on a chain specific designation. This collation is what must be relayed to the Consumer Chain via Cross Chain Validation.

The rate at which these validator set updates are relayed is a function of safety. At one extreme you could imagine the validator set being collated and relayed with every single epoch that is produced on the Provider Chain. This would ensure that the Consumer Chain has an absolutely accurate representation of validator weights at every potential state update within the Provider Chain. However, since delegations are subject to unbonding periods it is possible to approach state updates more conservatively. At another extreme, one may reason that if there is no active stake unbonding happening on the Provider Chain, it may be assumed that a validator set update will not be possible within a maximum duration of that unbonding period and so a validator set update can wait until just before that moment. Based on the possibility of instant redelegations this assumption may need to be further adjusted. 

The process of recording and relaying validator set updates within safe and correct periods is the focus of the spec and research at Informal Systems. We can assume that the design will be aware of the necessary cadence that validator set updates must take place to ensure safe operation. When there are no updates it may be necessary to simply acknowledge this with something like a heartbeat IBC packet.


### Reward Distribution

While Interchain Security may remove the necessity for a staking token to play a role in the token design of new blockchains, it can be assumed that there will be some sort of economic system in place to reward validators for producing blocks. To follow the default capabilities of the Cosmos SDK one would assume that there is a simple inflationary reward token attached to the production of blocks. This token may also be a governance token and the value may be implied by the ability to govern an otherwise useful and valuable network. The token could have many uses or reasons to exist, but being the responsibility of each blockchain to design its own purpose, we will assume some sort of reward token exists and is used to reward validators for the production of blocks.

In the current system rewards are pooled into the distribution module account. The distribution module imports the staking module to keep track of validator weights in order to calculate and distribute rewards on a per-delegator basis. For this to take place on the Consumer Chain it would be necessary for not only the validator set and validator voting power to be relayed with cross chain validation, but also all constituent delegators that compose each validator. This would result in an extremely large IBC packet which would make regular communication difficult if not impossible. Rather than relay all delegators, this information should stay on the Provider Chain within the Provider Chain staking module. 

Instead of distributing rewards on the Consumer Chain, the rewards collected for each block produced should be transferred to the Provider Chain at some interval. It may be a regular interval or a user initiated IBC packet, but would be similar to an ics-20 token transfer from the distribution module account of the Consumer Chain to the distribution module account of the Provider Chain. The difference between a standard ics-20 transfer being that the Provider Chain distribution module account needs to know which chain the rewards came from and over which period they were collected if this was not using V1 Full Validator Set Interchain Security. This information will be necessary to perform the Provider Chain distribution module's responsibility of distributed rewards to delegators and validators on a per chain basis informed by the Chain Registry.

### Evidence

In a single chain system, a validator may misbehave in various ways that result in the stake attached to that validator being slashed. This can occur automatically within the state machine or only after evidence of the misbehaviour has been collected and submitted to the chain. In a scenario where the validator set for a Consumer Chain is secured by stake that exists only on the Provider Chain, the evidence of misbehaviour needs to be submitted to the Provider Chain, where the tokens at stake are able to be slashed. Similar to a single chain system, this may take place automatically or only with the manual submission of evidence. If it were to occur automatically it would mean an outgoing IBC packet could be automatically generated but would still need to be manually relayed to the Provider Chain where punishment could follow. The more manual scenario would mean that the misbehaviour on the Consumer Chain results in evidence which is submitted directly to the Provider Chain, or submitted to the Consumer Chain which if successfully processed would result in an outgoing IBC packet containing the instruction to slash at the level of the Provider Chain.

Either way comes with the question of whether it should be the Provider Chain that is verifying the evidence of slashable events, or simply trusting the Consumer Chain's commands to slash validators at the level of the Provider Chain. While there may be the ability to both, trusting the Consumer Chain makes the job of the Provider Chain much easier. In that scenario the Provider Chain doesn't need to have knowledge of the logic contained within the Consumer Chain that may be necessary for determining whether a slashable offense took place. Trusting the Consumer Chain to enforce its own rules puts a lot of trust in the Consumer Chain to not abuse the position and submit a wave of slashable commands that may reduce a validator to nothing. However this is ultimately the responsibility of the validator to ensure that the slashable risk of validating on a Consumer Chain is worth the expected rewards. This includes determining whether the state machine of the Consumer Chain includes logic that the validator is comfortable with.

Limits on a Consumer Chain's ability to slash a Provider Chain validator may also be imposed at the Provider Chain level. For instance it is possible to make it a required parameter of a Consumer Chain within the Chain Registry to provide a maximum slashing rate. This would ensure that a validator knows even if they violate all the rules of a Consumer Chain, they can still limit the damage to their stake by some amount. It is also possible to enforce a parameter at the level of the Provider Chain that prevents a validator from joining any combination of Consumer Chain validator sets that results in a combined slashable rate over some threshold. For instance if there are 3 Consumer Chains that each have maximum slashable rates of 33% over X blocks, and the Provider Chain has a safety limit for validators that require them to stay below a total of 90% slashable events over X blocks, no validator would be allowed to join the Chain Registry for all 3 of those Consumer Chains. This parameter on the Provider Chain should be a governance parameter that can be adjusted by token staked voting to reflect the risk appetite of the actual validator set of the Provider Chain.

Consideration should be made to the incentives around submitting evidence in order to ensure that punishable offenses do not go unnoticed. This is a similar issue to Relayer Incentivization, where currently it costs some fee to relay IBC packets but there is no reward or payment possible as part of the core IBC logic. As a result IBC packets are currently relayed in an altruistic manner. It's important to ensure that for integrity of the operation of the child and Provider Chain that slashable offenses are always submitted. It may be that the slashable amounts of tokens are used as rewards for submitting the evidence (assuming it's possible to ensure it is not the culprit submitting their own evidence and regaining their stake). Maybe some flat fee to at least pay for the transactions are required, although this may become redundant with the addition of any Relayer incentivizations into core IBC.


## Consumer Chain Staking Module

With Interchain Security, the Consumer Chain receives updates of the Provider Chain's set of validator's voting power via cross chain validation IBC packets. These updates are used to populate the Consumer Chain's staking module. On the Consumer Chain, Tendermint consensus is running which asks the Consumer Chain staking module for the current list of validators and their voting power. This could work virtually the same as a traditional configuration without Interchain Security, as from the perspective of Tendermint the flow is the same (ABCI asks the staking module for this information).


## Consumer Chain Distribution Module

Delegators and validators on the Provider Chain would likely not risk their ATOMs in order to be included in the validator set on a Consumer Chain unless there was some incentive to do so. At a minimum, transaction fees gained from processing transactions could be considered as a possible incentive. There may also be incentives completely outside the state machine, like a service agreement that comes with regular payments through traditional means of money transmission. More likely it is expected that Consumer Chains include some type of block reward as seen in traditional validation schemes. This could be in a child-chain specific token used for gas, governance or some other use.

Regardless of how exactly a reward is calculated it is left up to the Consumer Chain to design a system that is attractive enough for Provider Chain validators to risk their staking token in order to be eligible. Once that reward is calculated it needs to be distributed back to the validators that have earned it.

In order to allow delegators to have a similar reward distribution as they currently experience, rewards should be regularly transferred to the Provider Chain where they can utilize the Provider Chain Distribution Module that tracks Cross Chain Validation client connections. These would be sent to the Provider Chain via a CCV IBC transaction similar to a traditional IBC Transfer packet but send to the Distribution Module Account with relevant information as to which validator and chain they were earned on and over what period of time. This information would allow the Provider Chain Distribution Module to distribute the rewards to the constituent delegators of each validator.

# Conclusion

Interchain Security is a new shared security primitive that has implications for the security and scalability of single blockchains like the Cosmos Hub as well as the potential to dramatically lower the barrier to running secure public blockchains for new applications. It could be thought of as a competitive configuration to sharding and put the Cosmos Hub on par with Eth 2.0 or Polkadot in terms of their security offerings to applications included in their environment. The design philosophy around the Cosmos ecosystem has always prioritized autonomy and sovereignty over guarantees around security—"Bring your own security" is a term that has been used in the past. The design of Interchain Security discussed here tries to incorporate the Cosmos design philosophy of autonomy and sovereignty with the offering of shared security. Even within the balance between those considerations there is a spectrum of possibilities and it might be the case that multiple versions and configurations are implemented in parallel in order to satisfy as many needs as possible.

# Further Reading

 * [Interchain Security is Coming to the Cosmos Hub](https://blog.cosmos.network/interchain-security-is-coming-to-the-cosmos-hub-f144c45fb035) - Billy Rennekamp
 * [Interchain Security Slidedeck from Cosmoverse Community Conference](https://docs.google.com/presentation/d/1XaPrbcNksnVdhZO1eyshyDDDQkA6buvKt90yxRF7sLs/edit?usp=sharing) - Billy Rennekamp
