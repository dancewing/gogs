(function(){
    window.md_merge_request_plugin = function(md) {
        function mrrule(state, silent) {
            var start = state.pos,
                regex = /([a-z0-9\.-_\/]+)?!(\d+)/,
                max   = state.posMax;

            if (state.src.charCodeAt(start) !== 0x21 || silent) {
                return false;
            }

            var matches = state.src.match(regex);

            if (! matches) {
                return false;
            }

            var id = matches[2];
            var title = matches[0];
            var board_name = state.env.host_url + '/' + state.env.board_url;
            if (matches[1]) {
                board_name =  state.env.host_url + '/' + matches[1];
                state.pending = state.pending.slice(0, -matches[1].length);
            }

            token = state.push('link_open', 'a');
            token.attrPush(['title', title]);
            token.attrPush(['href', board_name + '/merge_requests/' + id]);
            token.attrPush(['target', '_blank']);
            token.nesting = 1;

            token = state.push('text');
            token.content = title;
            token.nesting = 0;

            token = state.push('link_close', 'a');
            token.nesting = -1;

            state.pos = start + id.length +1;
            return true;
        }

        md.inline.ruler.push('mrrule', mrrule);
    };
}());
