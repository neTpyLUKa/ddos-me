pragma solidity >=0.7.0 <0.9.0;

    struct worker {
        uint256 id;
        int64 amount;
        address addr;
    }

    struct jobInfo {
        string host; // hostname of service to query
        int16 port;
        uint64 total; // total reward in wei
        uint64 start_time; // start_time, when service starts to process queires
        uint64 end_time; // not implemented yet
        uint64 registration_end_time; // time when registration ends
        mapping(address => worker) workers;
        uint64 total_registered;
        uint64 total_queries;
        // todo add timeout time?
    }

contract Child {
    mapping(address => uint256) private pendingReturns;

    uint256 constant public deposit_border = 10; // wei, 1 of 10^18

    constructor() {
    }

    function deposit() public payable {
        require(msg.value >= deposit_border); // TODO move this check to isAuthorized ?
        pendingReturns[msg.sender] += msg.value;
    }

    function withdraw() public returns (bool) {
        uint256 amount = pendingReturns[msg.sender];
        if (amount > 0) {
            pendingReturns[msg.sender] = 0;

            (bool result, ) = msg.sender.call{value:amount}("");
            if (!result) {
                pendingReturns[msg.sender] = amount;
                return false;
            }
        }
        return true;
    }

    function getDepositedValue(address addr) view public returns (uint256) {
        return pendingReturns[addr];
    }

    function isAuthorized(address addr) view public returns (bool) {
        return pendingReturns[addr] > deposit_border;
    }

    uint256 constant public queryReward = 1; // wei, 1 of 10^18

    event OnJobCreate(address creator); // remaining info will be pulled in read only mode

    mapping(address => jobInfo) public jobInfos;

    function createJob(string calldata host,
        int16 port,
        uint64 total,
        uint64 start_time,
        uint64 end_time,
        uint64 registration_end_time) public payable {
        require(isAuthorized(msg.sender));
        require(msg.value == queryReward * total); // TODO get total from queryReward and msg.value ?
        jobInfos[msg.sender].host = host;
        jobInfos[msg.sender].port = port;
        jobInfos[msg.sender].total = total;
        jobInfos[msg.sender].start_time = start_time;
        jobInfos[msg.sender].end_time = end_time;
        jobInfos[msg.sender].registration_end_time = registration_end_time;
        jobInfos[msg.sender].total_registered = 0;
        jobInfos[msg.sender].total_queries = 0;

        emit OnJobCreate(msg.sender);
    }

    function register(address jobAddress) public {
        require(isAuthorized(msg.sender));
        jobInfo storage info = jobInfos[jobAddress];
        require(block.timestamp < info.registration_end_time);

        worker storage w = info.workers[msg.sender];
        w.id = 0;
        //        require(w.amount == 0);

        w.addr = msg.sender;
        w.amount = -1;

        info.total_registered += 1;
    }

    function isRegistrationEnded(address jobAddress) public view returns (bool) {
        return jobInfos[jobAddress].registration_end_time >= block.timestamp;
    }

    function isRegistered(address jobAddress) public view returns (bool) {
        return jobInfos[jobAddress].workers[jobAddress].amount != 0;
    }

    function submitResult(uint256 id, int64 amount, address jobAddress) public { // TODO accept signature
        // TODO нужно разрешить сабмитить только до некоторого timeout, проверка подписи там
        require(jobInfos[jobAddress].total != 0); // TODO make general function which checks if job exists
        jobInfo storage info = jobInfos[jobAddress];
        worker storage w = info.workers[msg.sender];
        require((w.id & id) == 0);

        uint64 value = uint64(amount); // TODO check if amount > 0
        w.amount += amount;
        w.id |= id;
        w.addr = msg.sender; // needed for first submit

        int64 border = int64(info.total / info.total_registered); // TODO store this in jobInfo
        if (w.amount > border) {
            value -= uint64(w.amount - border); // TODO fix types
            w.amount = border;
        }
        info.total_queries += value;
    }

    // TODO интеграционные тесты

    function claimReward(address jobAddress) public {
        // TODO разрешить забирать награду только после сабмита некоторого
        // TODO защита от двойных списаний
        jobInfo storage info = jobInfos[jobAddress];

        worker storage w = info.workers[msg.sender];
        require(w.amount > 0);
        require(info.total_queries >= w.amount); // for debug

        uint64 value = uint64(uint128(uint64(w.amount)) * uint128(info.total) / uint128(info.total_queries));

        info.total -= value;
        info.total_queries -= uint64(w.amount);
        w.amount = 0;

        msg.sender.call{value:value}(""); // TODO handle return result error

        // require(block.timestamp > info.end_time);

    }
}