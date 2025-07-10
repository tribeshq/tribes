// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.21;

import {Proof} from "vlayer-1.2.0/src/Proof.sol";
import {Prover} from "vlayer-1.2.0/src/Prover.sol";
import {Web, WebProof, WebProofLib, WebLib} from "vlayer-1.2.0/src/WebProof.sol";

contract WebProofXProver is Prover {
    using WebProofLib for WebProof;
    using WebLib for Web;

    string public constant DATA_URL =
        "https://api.x.com/1.1/account/settings.json?include_ext_sharing_audiospaces_listening_data_with_followers=true&include_mention_filter=true&include_nsfw_user_flag=true&include_nsfw_admin_flag=true&include_ranked_timeline=true&include_alt_text_compose=true&ext=ssoConnections&include_country_code=true&include_ext_dm_nsfw_media_filter=true";

    function main(WebProof calldata webProof, address account, address application)
        public
        view
        returns (Proof memory, string memory, address, address)
    {
        Web memory web = webProof.verify(DATA_URL);

        string memory screenName = web.jsonGetString("screen_name");

        return (proof(), screenName, account, application);
    }
}
