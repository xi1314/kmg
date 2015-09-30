/*
 import(
 "github.com/bronze1man/kmg/kmgView/kmgWeb/jquery"
 )
 */

//初始化webuploader
$(function(){
    $(".upload-image").each(function(){
        var comp = $(this);
        var targetUrl = comp.attr("uploadUrl");
        var fileInput = comp.find(".upload-image-btn");
        var progress = comp.find(".progress")[0];
        var messageDiv = comp.find('.upload-message');
        var image = comp.find('img');
        var delteBtn = comp.find(".btn-delete-img");

        $(delteBtn).on("click",function(){
            $(image).addClass("hide");
            comp.find('.upload-set-input').val("");
            messageDiv.text("");
            $(this).hide();
        });

        var uploader = WebUploader.create({
            auto: true,
            swf: getResourceUrlPrefix()+'/WebUploader.swf',
            server: targetUrl,
            pick: {
                id: fileInput,
                multiple: false
            },
            thumb: false,
            resize: false,
            compress: false,
            formData:{
                "ProcessorName":  comp.attr('imageProcessorName')
            },
            headers: {'X-A':'A'} //ie7下设置header,才会发referer
        });
        uploader.on( 'uploadProgress', function( file, percentage ) {

            $(progress).show();
            $(delteBtn).hide();
            $(image).addClass("hide");

            $(progress).css({width:(percentage * 90) + '%'});
        });
        uploader.on( 'uploadSuccess', function( file,ret ) {
            $(progress).css({width:'100%'});

            $(image).attr("src",ret.data).removeClass("hide");
            comp.find('.upload-set-input').val(ret.data);

            $(delteBtn).show();
            $(messageDiv).text("上传成功");
            $(progress).hide();

        });
        uploader.on( 'uploadError', function( file) {
            $(messageDiv).text(file.name+"上传失败");
            $(progress).hide();
        });
    })
});